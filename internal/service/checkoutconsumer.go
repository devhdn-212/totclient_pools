package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/devhdn-212/totclient_pools/internal/repository"
	"github.com/devhdn-212/totclient_pools/internal/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// rejectedCheckoutTTL is deliberately short (2 jam) — this is a
// troubleshooting/reconciliation aid for "player saw accepted but bet was
// never persisted" cases (see DOKUMENTASI.md § 3.11/§ 9.3), not a durable
// record. Postgres is never touched for a balance-rejected event, so there
// is nothing more permanent to fall back to once this expires.
const rejectedCheckoutTTL = 2 * time.Hour

// rejectedCheckoutRecord is what gets cached when CheckBalance finds the
// player's live balance no longer covers this chunk — the full event plus
// why, so whoever investigates a "player says bet didn't go through" report
// doesn't have to reconstruct it from log lines alone.
type rejectedCheckoutRecord struct {
	Reason     string                      `json:"reason"`
	RejectedAt time.Time                   `json:"rejected_at"`
	Event      domain.CheckoutKafkaMessage `json:"event"`
}

func rejectedCheckoutCacheKey(idcomp string, idtrxkeluaran, playerinvoice int) string {
	return "client:rejected:" + strings.ToLower(idcomp) + ":" + strconv.Itoa(idtrxkeluaran) + ":playerinvoice:" + strconv.Itoa(playerinvoice)
}

// redisAgenTrxkeluaranmemberPrefix / trxkeluaranmemberByUserCacheKey mirror
// the exact cache keys totclient_api's HTTP-facing trxkeluaranmemberService
// reads (see that repo's internal/service/trxkeluaranmember.go) — this
// worker is what now actually writes the rows behind those caches, so it's
// the one that has to bust them, even though it never reads them itself.
const (
	redisAgenTrxkeluaranmemberPrefix = "agen:trxkeluaranmember"
	redisTrxkeluaranmemberByUser     = "client:trxkeluaranmember"
)

// redisRiwayatTransaksiPrefix / redisRiwayatDetailTransaksiPrefix mirror the
// exact cache keys totclient_api's riwayatTransaksiService.Fetch reads
// (Redisriwayattransaksi / Redisriwayatdetailtransaksi there, see that
// repo's internal/service/riwayattransaksi.go) — same mirroring rationale
// as redisAgenTrxkeluaranmemberPrefix above: this worker is what actually
// persists the bet a checkout represents (trxkeluarandetail/
// trxkeluaranmember), so it's the one that has to bust these two too, even
// though it never reads them itself.
const (
	redisRiwayatTransaksiPrefix       = "client:riwayattransaksi"
	redisRiwayatDetailTransaksiPrefix = "client:riwayattransaksi:detail"
)

func trxkeluaranmemberByUserCacheKey(idcomp, idcomppasaran, username string) string {
	return redisTrxkeluaranmemberByUser + ":v4:" + strings.ToLower(idcomp) + ":" + idcomppasaran + ":" + username
}

// CheckoutConsumer is the Kafka consumer side of the checkout flow: it reads
// CheckoutKafkaMessage events totclient_api publishes once balance/limit
// validation passes, debits the wallet, and — only on a successful debit —
// persists the bet (trxkeluarandetail/trxkeluaranmember) and refreshes the
// period totals. See domain/checkout.go for the full contract.
type CheckoutConsumer struct {
	db                *pgxpool.Pool
	reader            *kafka.Reader
	memberinfoService domain.MemberinfoService
}

func NewCheckoutConsumer(db *pgxpool.Pool, reader *kafka.Reader, memberinfoService domain.MemberinfoService) *CheckoutConsumer {
	return &CheckoutConsumer{
		db:                db,
		reader:            reader,
		memberinfoService: memberinfoService,
	}
}

// Run blocks, consuming checkout events until ctx is canceled — meant to be
// the process's main loop (see main.go), not something spun off into a
// background goroutine and forgotten.
// fetchRetryBackoff is how long Run waits after a failed FetchMessage
// before trying again — without it, a persistent broker outage/DNS failure
// turns this into a tight busy-loop hammering the broker (and spamming the
// Telegram alert hook, see below) on every single iteration instead of a
// bounded rate.
const fetchRetryBackoff = 5 * time.Second

func (c *CheckoutConsumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			// The error detail goes IN the message string, not as a zap
			// field — connection.Log.Error's Telegram hook only forwards
			// entry.Message (see DOKUMENTASI.md § 3.7/3.8), so a plain
			// zap.Error(err) here would make every alert say the same
			// useless "FetchMessage failed" with no indication of why
			// (broker unreachable, DNS failure, topic missing, etc.).
			connection.Log.Error("checkout consumer: FetchMessage failed: " + err.Error())
			select {
			case <-ctx.Done():
				return
			case <-time.After(fetchRetryBackoff):
			}
			continue
		}

		c.handle(ctx, msg.Value)

		// Committed unconditionally (success or a logged wallet/persist
		// failure) — there is no retry/dead-letter path here, matching how
		// a wallet rejection is handled below: it's a deliberate,
		// print-and-move-on decision, not an oversight.
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			connection.Log.Error("checkout consumer: CommitMessages failed: " + err.Error())
		}
	}
}

func (c *CheckoutConsumer) handle(ctx context.Context, raw []byte) {
	var evt domain.CheckoutKafkaMessage
	if err := json.Unmarshal(raw, &evt); err != nil {
		connection.Log.Error("checkout consumer: invalid message payload", zap.Error(err))
		return
	}

	fmt.Printf("Received message from %s channel with invoice-player: %d of user: %s\n", evt.Idcompany, evt.Playerinvoice, evt.Username)

	// Balance is realtime — re-check it right before debiting instead of
	// trusting the balance totclient_api validated against when this event
	// was published (could be stale by the time this message is actually
	// consumed, e.g. Kafka/consumer lag, or the player spent balance
	// elsewhere in the meantime).
	balance, err := c.memberinfoService.CheckBalance(ctx, evt.Idcompany, evt.Token)
	if err != nil {
		connection.Log.Error("checkout consumer: balance check failed",
			zap.String("playerinvoice", evt.PlayerinvoiceStr), zap.Error(err))
		return
	}
	if balance.LessThan(evt.Totalbayar) {
		fmt.Printf("pasaran : %s, invoice : %d\n", evt.Pasaran, evt.IDtrxkeluaran)
		fmt.Println("invoice ini balance tidak cukup")
		fmt.Printf("%s : failed execute transaction for invoiceid: %d\n", evt.Username, evt.Playerinvoice)

		record := rejectedCheckoutRecord{
			Reason:     "invoice ini balance tidak cukup",
			RejectedAt: util.GetNowJakarta(),
			Event:      evt,
		}
		key := rejectedCheckoutCacheKey(evt.Idcompany, evt.IDtrxkeluaran, evt.Playerinvoice)
		if err := connection.SetRedis(key, record, rejectedCheckoutTTL); err != nil {
			connection.Log.Error("checkout consumer: failed to cache rejected checkout",
				zap.String("playerinvoice", evt.PlayerinvoiceStr), zap.Error(err))
		}
		return
	}

	_, err = c.memberinfoService.SubmitTransaction(ctx, domain.MothershipTransaction{
		Idcompany:     evt.Idcompany,
		Invoice:       evt.Invoice,
		Pasaran:       evt.Pasaran,
		Playerinvoice: evt.PlayerinvoiceStr,
		Username:      evt.Username,
		Credit:        decimal.Zero,
		Debit:         evt.Totalbayar,
	})
	if err != nil {
		// Wallet rejected (or the call itself failed) — per product
		// decision this bet just never happened as far as
		// trxkeluarandetail/trxkeluaranmember are concerned. No retry, no
		// dead-letter, nothing persisted; only a log line survives.
		fmt.Println("Checkout consumer - wallet GAGAL, playerinvoice:", evt.PlayerinvoiceStr, "username:", evt.Username, "error:", err)
		return
	}

	if err := c.persist(ctx, evt); err != nil {
		connection.Log.Error("checkout consumer: failed to persist checkout after wallet debit succeeded",
			zap.String("playerinvoice", evt.PlayerinvoiceStr), zap.Error(err))
		return
	}

	fmt.Printf("%s : finish execute transaction for invoiceid: %d\n", evt.Username, evt.Playerinvoice)
}

// persist inserts one trxkeluarandetail row per item plus the
// trxkeluaranmember summary row, all-or-nothing in one tx — mirrors exactly
// what totclient_api's checkoutService.Submit used to do inline before this
// was split across Kafka, just fed from the event instead of a live
// pasaran-rate calculation (that already happened in totclient_api).
func (c *CheckoutConsumer) persist(ctx context.Context, evt domain.CheckoutKafkaMessage) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	detailRepo := repository.NewTrxkeluarandetailRepository(txExec)
	memberRepo := repository.NewTrxkeluaranmemberRepository(txExec)

	compileItems := make([]compileDetailItem, 0, len(evt.Items))
	for _, item := range evt.Items {
		detail := &domain.Trxkeluarandetail{
			ID:                   item.DetailID,
			IDtrxkeluaran:        evt.IDtrxkeluaran,
			IDcomp:               evt.Idcompany,
			Datekeluarandetail:   sql.NullTime{Valid: true, Time: evt.Datekeluarandetail},
			Ipaddress:            evt.Ipaddress,
			Username:             evt.Username,
			Typegame:             item.Typegame,
			Nomortogel:           item.Nomor,
			Posisitogel:          item.Posisitogel,
			Bet:                  item.Bet,
			Diskon:               item.DiscRate,
			Win:                  item.Win,
			Kei:                  item.KeiRate,
			Browsertogel:         evt.Devicemember,
			Devicetogel:          evt.Devicemember,
			Statuskeluarandetail: "RUNNING",
			Betround:             evt.Betround,
			Playerinvoice:        evt.Playerinvoice,
			Created:              evt.Username,
			CreatedAt:            sql.NullTime{Valid: true, Time: evt.Datekeluarandetail},
		}
		if err := detailRepo.Save(ctx, detail, evt.Idcompany); err != nil {
			return fmt.Errorf("insert trxkeluarandetail: %w", err)
		}
		compileItems = append(compileItems, compileDetailItem{
			Idtrxdetail:   detail.ID,
			Playerinvoice: evt.PlayerinvoiceStr,
			Username:      evt.Username,
			Nomor:         item.Nomor,
			Type:          item.Typegame,
			Posisi:        item.Posisitogel,
			Bet:           item.Bet,
			Payout:        item.Payout,
			Win:           item.Win,
		})
	}

	member := &domain.Trxkeluaranmember{
		ID:            evt.MemberID,
		IDtrxkeluaran: evt.IDtrxkeluaran,
		IDcomp:        evt.Idcompany,
		Username:      evt.Username,
		Totalbet:      evt.Totalbet,
		Totalbayar:    evt.Totalbayar,
		Totaldiscount: evt.Totaldiscount,
		Totalkei:      evt.Totalkei,
		Totalpair:     evt.Totalpair,
		Betround:      evt.Betround,
		Playerinvoice: evt.Playerinvoice,
		Status:        "pending",
		Created:       evt.Username,
		CreatedAt:     sql.NullTime{Valid: true, Time: evt.Datekeluarandetail},
	}
	if err := memberRepo.Save(ctx, member, evt.Idcompany); err != nil {
		return fmt.Errorf("insert trxkeluaranmember: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Everything below is post-commit cache maintenance — best-effort,
	// fire-and-forget, never allowed to turn a persisted bet into a
	// reported failure.
	go connection.DeleteRedis(trxkeluaranmemberByUserCacheKey(evt.Idcompany, evt.PasaranIdcomp, evt.Username))
	go connection.DeleteRedisByPrefix(redisAgenTrxkeluaranmemberPrefix + ":" + strings.ToLower(evt.Idcompany))
	// Checkout just added a new bet row for this user/pasaran/period —
	// scoped down to username (unlike the savenomorkeluaran/revisi
	// invalidation on the agen side, which busts every player's cache
	// because a draw settling changes everyone's status at once) since a
	// new bet here only ever changes what THIS player's own riwayat looks
	// like.
	//
	// Level-1 (period list) is deleted by EXACT key, not prefix: the cache
	// key riwayatTransaksiService.Fetch writes for this case
	// ("client:riwayattransaksi:<idcomp>:<pasaran_idcomp>:<username>") has
	// no segment after username, so DeleteRedisByPrefix's "prefix:*" scan
	// pattern would never match it (that trailing ":*" requires something
	// AFTER the prefix) — using it here silently never deleted anything,
	// which is exactly why the Transaksi list kept showing stale data after
	// a checkout. Level-2 (detail) genuinely does have a segment after this
	// same prefix (":<idtrxkeluaran>"), so prefix-delete stays correct there.
	go connection.DeleteRedis(redisRiwayatTransaksiPrefix + ":" + strings.ToLower(evt.Idcompany) + ":" + evt.PasaranIdcomp + ":" + evt.Username)
	go connection.DeleteRedisByPrefix(redisRiwayatDetailTransaksiPrefix + ":" + strings.ToLower(evt.Idcompany) + ":" + evt.PasaranIdcomp + ":" + evt.Username)
	if len(compileItems) > 0 {
		go cacheCompileData(evt.Idcompany, evt.IDtrxkeluaran, evt.PlayerinvoiceStr, compileItems)
	}
	MarkTotalsDirty(evt.Idcompany, evt.IDtrxkeluaran)

	return nil
}
