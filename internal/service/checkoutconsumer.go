package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

// errInvoiceAlreadyVoided marks the specific persist() failure where an
// admin's manual refund won the race and voided this invoice before persist()
// could complete (see MarkCompleted in internal/repository/memberinvoice.go
// and DOKUMENTASI.md § 9.9). Unlike other persist() failures, retrying this
// one is pointless (the row will never become Requested again) and it is not
// a "money taken, data missing" case needing reconciliation - the money
// already went back to the player, correctly.
var errInvoiceAlreadyVoided = errors.New("member invoice already voided by admin refund - aborting persist")

// CheckoutConsumer is the Kafka consumer side of the checkout flow: it reads
// CheckoutKafkaMessage events totclient_api publishes once balance/limit
// validation passes, debits the wallet, and — only on a successful debit —
// persists the bet (trxkeluarandetail/trxkeluaranmember) and refreshes the
// period totals. See domain/checkout.go for the full contract.
type CheckoutConsumer struct {
	db                *pgxpool.Pool
	reader            *kafka.Reader
	memberinfoService domain.MemberinfoService
	memberInvoiceRepo domain.MemberInvoiceRepository
}

func NewCheckoutConsumer(db *pgxpool.Pool, reader *kafka.Reader, memberinfoService domain.MemberinfoService, memberInvoiceRepo domain.MemberInvoiceRepository) *CheckoutConsumer {
	return &CheckoutConsumer{
		db:                db,
		reader:            reader,
		memberinfoService: memberinfoService,
		memberInvoiceRepo: memberInvoiceRepo,
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

// processTimeout bounds the "finish this one message no matter what" context
// below - long enough for a normal debit+persist, short enough that a truly
// stuck DB/wallet call doesn't hang shutdown forever (orchestrators kill the
// process outright after their own grace period anyway).
const processTimeout = 30 * time.Second

// persistMaxAttempts/persistRetryBackoff - dikonfirmasi user: cuma buat
// persist() (nyimpen trxkeluarandetail/trxkeluaranmember) SETELAH debit
// wallet sudah sukses, bukan buat SubmitTransaction itu sendiri (retry debit
// punya risiko beda - kalau percobaan pertama sebenarnya sukses tapi
// response-nya hilang, retry bisa kirim debit dobel; persist() tidak punya
// risiko itu karena baru dipanggil SETELAH debit dipastikan sukses). 3x
// percobaan @1 detik - total tambahan waktu ~2 detik di kasus terburuk,
// jauh di bawah processTimeout (30 detik) di atas.
const (
	persistMaxAttempts  = 3
	persistRetryBackoff = 1 * time.Second
)

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

		// Sengaja pakai context BARU (bukan ctx yang dibatalkan pas shutdown/
		// redeploy) buat memproses SATU message yang sudah kadung ke-fetch -
		// dikonfirmasi user ini penyebab nyata sebagian baris "Pending Bet
		// Request" (DOKUMENTASI.md § 9.8/§13): kalau SIGTERM datang PERSIS di
		// tengah handle()/persist(), pola lama (pakai ctx yang sama buat
		// FetchMessage DAN proses+commit) bisa motong transaksi persist di
		// tengah jalan walau debit wallet-nya udah kadung sukses - uang
		// kepotong, data hilang, padahal seharusnya bisa dihindari sepenuhnya
		// kalau restart-nya nunggu message yang lagi diproses selesai dulu.
		// Message yang sudah di-fetch sekarang SELALU diselesaikan sampai
		// tuntas (debit+persist+commit offset), baru loop mengecek ctx.Done()
		// SEBELUM fetch message berikutnya - bukan berhenti di tengah message
		// yang sedang diproses.
		processCtx, cancel := context.WithTimeout(context.Background(), processTimeout)
		c.handle(processCtx, msg.Value)

		// Committed unconditionally (success or a logged wallet/persist
		// failure) — there is no retry/dead-letter path here, matching how
		// a wallet rejection is handled below: it's a deliberate,
		// print-and-move-on decision, not an oversight.
		if err := c.reader.CommitMessages(processCtx, msg); err != nil {
			connection.Log.Error("checkout consumer: CommitMessages failed: " + err.Error())
		}
		cancel()

		if ctx.Err() != nil {
			return
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

	// dateTransaction anchors this checkout's row in public.tbl_trx_member_invoice
	// (PK/partition column) - captured once here so InsertRequested/MarkVoid/
	// MarkCompleted below all address the exact same row.
	now := util.GetNowJakarta()
	dateTransaction := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

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
		//
		// Deliberately NO row written to tbl_trx_member_invoice on this path
		// - dikoreksi user: Requested HARUS berarti "dana user sudah
		// kepotong", titik. Kalau baris ini ditulis di sini (sebelum tahu
		// debit-nya sukses atau tidak), baris "debit gagal" numpuk bareng
		// baris "debit sukses tapi persist gagal" di panel Pending Bet
		// Request, dan admin refund manual bisa KELIRU ngasih kredit ke
		// invoice yang uangnya belum pernah kepotong sama sekali - itu sama
		// aja ngasih uang gratis ke player. InsertRequested sengaja dipindah
		// ke SETELAH SubmitTransaction sukses (di bawah), bukan sebelum -
		// lihat DOKUMENTASI.md § 9.8.
		fmt.Println("Checkout consumer - wallet GAGAL, playerinvoice:", evt.PlayerinvoiceStr, "username:", evt.Username, "error:", err)
		return
	}

	// Requested row written HERE, SETELAH debit sukses (bukan sebelumnya) -
	// deliberately no I/O between the successful SubmitTransaction above and
	// this insert, so the window where a debit succeeded but no row exists
	// yet is as close to zero as this code can make it. This ordering
	// guarantees the OPPOSITE property from before: any row that DOES exist
	// here is guaranteed to represent "wallet was actually debited" - never
	// a debit attempt that failed/never happened. See DOKUMENTASI.md § 9.8.
	if err := c.memberInvoiceRepo.InsertRequested(ctx, &domain.MemberInvoice{
		AgentCode:       evt.Idcompany,
		InvoiceID:       evt.Playerinvoice,
		Username:        evt.Username,
		PlayerToken:     evt.Token,
		DebitAmount:     evt.Totalbayar,
		DateTransaction: dateTransaction,
	}); err != nil {
		connection.Log.Error("checkout consumer: failed to insert member invoice (Requested) after debit succeeded",
			zap.String("playerinvoice", evt.PlayerinvoiceStr), zap.Error(err))
		// Deliberately NOT returning here - the debit already succeeded, so
		// persist() below still has to be attempted regardless of whether
		// this tracking insert worked. Skipping persist() on top of a failed
		// tracking insert would double the damage of one failure (lost bet
		// data AND no reconciliation record) instead of just losing the
		// reconciliation record. persist()'s MarkCompleted no-ops harmlessly
		// if this row never existed (0 rows matched, not an error).
	}

	// persist() retried a few times before giving up - dikonfirmasi user:
	// debit sudah kadung sukses, jadi kalau gagal nyimpen gara-gara gangguan
	// SESAAT (DB blip, dsb - bukan lagi soal restart/shutdown, itu sudah
	// ditutup lewat processCtx di Run()), coba lagi dulu beberapa kali
	// sebelum benar-benar nyerah jadi baris "Requested". persist() buka
	// transaksi baru sendiri tiap dipanggil, jadi aman dipanggil ulang -
	// tidak ada state parsial yang perlu dibersihkan di antara percobaan.
	var persistErr error
	for attempt := 1; attempt <= persistMaxAttempts; attempt++ {
		persistErr = c.persist(ctx, evt, dateTransaction)
		if persistErr == nil {
			break
		}
		if errors.Is(persistErr, errInvoiceAlreadyVoided) {
			// Retrying is pointless - the row is Void now and will never go
			// back to Requested, so every retry would fail the exact same
			// way. Stop immediately instead of burning persistMaxAttempts
			// tries and backoff sleeps on a foregone conclusion.
			break
		}
		if attempt < persistMaxAttempts {
			connection.Log.Warn(fmt.Sprintf("checkout consumer: persist gagal (percobaan %d/%d), retry dalam %s", attempt, persistMaxAttempts, persistRetryBackoff),
				zap.String("playerinvoice", evt.PlayerinvoiceStr), zap.Error(persistErr))
			time.Sleep(persistRetryBackoff)
		}
	}
	if persistErr != nil {
		if errors.Is(persistErr, errInvoiceAlreadyVoided) {
			// Not a failure needing reconciliation - an admin already
			// resolved this invoice (refund) before persist() got to it.
			// The bet is correctly NOT persisted; the money is correctly
			// back in the player's wallet. See DOKUMENTASI.md § 9.9.
			connection.Log.Info("checkout consumer: invoice sudah di-void (refund manual admin) sebelum sempat dipersist - bet sengaja tidak disimpan",
				zap.String("playerinvoice", evt.PlayerinvoiceStr))
			return
		}
		connection.Log.Error(fmt.Sprintf("checkout consumer: failed to persist checkout after wallet debit succeeded (sudah dicoba %d kali)", persistMaxAttempts),
			zap.String("playerinvoice", evt.PlayerinvoiceStr), zap.Error(persistErr))
		// Deliberately left at Requested (not voided) here - the debit
		// DID succeed, so this is exactly the "money taken, data missing"
		// case the member-invoice table exists to catch. Do not touch its
		// status; totagen_api's reconciliation is what resolves this.
		return
	}

	fmt.Printf("%s : finish execute transaction for invoiceid: %d\n", evt.Username, evt.Playerinvoice)
}

// persist inserts one trxkeluarandetail row per item plus the
// trxkeluaranmember summary row, all-or-nothing in one tx — mirrors exactly
// what totclient_api's checkoutService.Submit used to do inline before this
// was split across Kafka, just fed from the event instead of a live
// pasaran-rate calculation (that already happened in totclient_api).
func (c *CheckoutConsumer) persist(ctx context.Context, evt domain.CheckoutKafkaMessage, dateTransaction time.Time) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	detailRepo := repository.NewTrxkeluarandetailRepository(txExec)
	memberRepo := repository.NewTrxkeluaranmemberRepository(txExec)
	memberInvoiceRepo := repository.NewMemberInvoiceRepository(txExec)

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

	// Completed in the SAME tx as the two inserts above - atomic: either all
	// three land together, or (on any earlier error/rollback) the Requested
	// row from handle() stays exactly as it was, still correctly readable by
	// totagen_api's reconciliation.
	claimed, err := memberInvoiceRepo.MarkCompleted(ctx, evt.Idcompany, evt.Playerinvoice, dateTransaction)
	if err != nil {
		return fmt.Errorf("mark member invoice completed: %w", err)
	}
	if !claimed {
		// An admin's manual refund (totagen_api ClaimForRefund) already won
		// the race on this exact row and flipped it to Void - money has
		// already gone back to the player's wallet. Abort here (defer
		// tx.Rollback above undoes the trxkeluarandetail/trxkeluaranmember
		// inserts too) so this bet never becomes valid data on top of an
		// already-issued refund. See DOKUMENTASI.md § 9.9.
		return errInvoiceAlreadyVoided
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
