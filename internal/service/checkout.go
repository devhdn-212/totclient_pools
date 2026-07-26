package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/repository"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type checkoutService struct {
	db                *pgxpool.Pool
	trxkeluaranRepo   domain.TrxkeluaranRepository
	pasaranService    domain.PasaranService
	memberinfoService domain.MemberinfoService
}

func NewCheckoutService(
	db *pgxpool.Pool,
	trxkeluaranRepo domain.TrxkeluaranRepository,
	pasaranService domain.PasaranService,
	memberinfoService domain.MemberinfoService,
) domain.CheckoutService {
	return &checkoutService{
		db:                db,
		trxkeluaranRepo:   trxkeluaranRepo,
		pasaranService:    pasaranService,
		memberinfoService: memberinfoService,
	}
}

// resolveLimits maps a typegame to its (limittotal, limitglobal) ceiling
// from the pasaran config. Field names are wildly inconsistent across bet
// types (same issue as disc/kei/win elsewhere in this codebase), so each
// type is spelled out explicitly.
func resolveLimits(p dto.PasaranData, typegame string) (limittotal, limitglobal decimal.Decimal) {
	switch typegame {
	case "4D":
		return p.AngkaLimittotal4d, p.AngkaLimitbuang4d
	case "3D":
		return p.AngkaLimittotal3d, p.AngkaLimitbuang3d
	case "3DD":
		return p.AngkaLimittotal3dd, p.AngkaLimitbuang3dd
	case "2D":
		return p.AngkaLimittotal2d, p.AngkaLimitbuang2d
	case "2DD":
		return p.AngkaLimittotal2dd, p.AngkaLimitbuang2dd
	case "2DT":
		return p.AngkaLimittotal2dt, p.AngkaLimitbuang2dt
	case "COLOK_BEBAS":
		return p.CbLimitotal, p.CbLimitbuang
	case "COLOK_MACAU":
		return p.CmacauLimittotal, p.CmacauLimitbuang
	case "COLOK_NAGA":
		return p.CnagaLimittotal, p.CnagaLimitbuang
	case "COLOK_JITU":
		return p.CjituLimitotal, p.CjituLimitbuang
	case "50_50_UMUM":
		return p.Umum5050Limittotal, p.Umum5050Limitbuang
	case "50_50_SPECIAL":
		return p.Special5050Limittotal, p.Special5050Limitbuang
	case "50_50_KOMBINASI":
		return p.Kombinasi5050Limittotal, p.Kombinasi5050Limitbuang
	case "MACAU_KOMBINASI":
		return p.MacaukombinasiLimittotal, p.MacaukombinasiLimitbuang
	case "DASAR":
		return p.DasarLimittotal, p.DasarLimitbuang
	case "SHIO":
		return p.ShioLimittotal, p.ShioLimitbuang
	default:
		return decimal.Zero, decimal.Zero
	}
}

// stripAsterisks removes the "*" wildcard markers used in 3DD/2DD/2DT
// numbers (e.g. "869*", "*08*") — typegame already captures that
// distinction, so the stored nomortogel doesn't need them.
func stripAsterisks(nomor string) string {
	return strings.ReplaceAll(nomor, "*", "")
}

// isPasaranAcceptingBets re-derives whether the pasaran is open for betting
// RIGHT NOW instead of trusting p.Status — that field is baked into
// pasaranService.FindID's 24h cache at whatever moment the cache was last
// populated, so a market that has since passed its closing time could still
// read back as "ONLINE" for up to 24h. JadwalOpen/JadwalTutup are absolute
// timestamps for today's session (safe to cache — they don't change during
// the day), so comparing them against the live current time here is always
// accurate no matter how stale the surrounding cached PasaranData is.
//
// Falls back to the cached Status only when no (parseable) schedule exists
// for this pasaran/day — some markets run without a jam-buka/tutup window
// and rely solely on the admin-set status flag.
func isPasaranAcceptingBets(p dto.PasaranData, now time.Time) bool {
	openAt, errOpen := time.ParseInLocation("2006-01-02 15:04:05", p.JadwalOpen, util.LocJakarta)
	closeAt, errClose := time.ParseInLocation("2006-01-02 15:04:05", p.JadwalTutup, util.LocJakarta)
	if errOpen != nil || errClose != nil {
		return p.Status == "ONLINE"
	}
	return !now.Before(openAt) && now.Before(closeAt)
}

func (s *checkoutService) Submit(ctx context.Context, req dto.CheckoutRequest, ipaddress string) (dto.CheckoutResponse, error) {
	// idcompany/codecomppasaran are stored uppercase — /api/servicetoken and
	// /api/serviceinit already uppercase these before calling upstream, but
	// this endpoint gets the raw client payload (e.g. agent=nuk from the
	// launch URL), so it has to normalize the same way here.
	idcomp := strings.ToUpper(req.Company)
	pasaranCode := strings.ToUpper(req.PasaranCode)

	memberinfo, err := s.memberinfoService.CheckToken(ctx, dto.MemberinfoResponse{
		Agen:   idcomp,
		Market: pasaranCode,
		Token:  req.Token,
	})
	if err != nil {
		return dto.CheckoutResponse{}, err
	}
	if req.Total.GreaterThan(decimal.NewFromInt(memberinfo.Balance)) {
		return dto.CheckoutResponse{}, domain.ErrInsufficientBalance
	}
	username := memberinfo.Username

	// Cached (24h TTL + singleflight) — every player betting on the same
	// pasaran shares this instead of each checkout hitting the DB for the
	// same rate config.
	pasaran, err := s.pasaranService.FindID(ctx, idcomp, pasaranCode)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}

	now := util.GetNowJakarta()
	if !isPasaranAcceptingBets(pasaran, now) {
		return dto.CheckoutResponse{}, domain.ErrPasaranOffline
	}

	trxkeluaran, err := s.trxkeluaranRepo.FindByID(ctx, idcomp, req.PasaranIdcomp)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}
	if trxkeluaran.ID == 0 {
		return dto.CheckoutResponse{}, fmt.Errorf("trxkeluaran %w", util.ErrNotFound)
	}
	idtrxkeluaran := trxkeluaran.ID

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	detailRepo := repository.NewTrxkeluarandetailRepository(txExec)
	memberRepo := repository.NewTrxkeluaranmemberRepository(txExec)
	trxkeluaranTxRepo := repository.NewTrxkeluaranRepository(txExec)

	invoiceCounterKey := "tbl_trx_keluarantogel_detail_" + strings.ToLower(req.Company) + "_" + strconv.Itoa(idtrxkeluaran)
	invoiceCounter, err := util.GetNextCounterManualTx(ctx, tx, invoiceCounterKey)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}
	playerinvoiceStr := fmt.Sprintf("%d%d", idtrxkeluaran, invoiceCounter)
	playerinvoice, err := strconv.Atoi(playerinvoiceStr)
	if err != nil {
		return dto.CheckoutResponse{}, fmt.Errorf("error building playerinvoice: %w", err)
	}

	// How many playerinvoices this username has made within THIS invoice
	// game — resets back to 1 whenever idtrxkeluaran changes.
	betroundKey := "betround_" + strings.ToLower(req.Company) + "_" + strconv.Itoa(idtrxkeluaran) + "_" + username
	betroundCounter, err := util.GetNextCounterManualTx(ctx, tx, betroundKey)
	if err != nil {
		return dto.CheckoutResponse{}, err
	}
	betround := int(betroundCounter)

	// now was captured earlier (right before isPasaranAcceptingBets) — reused
	// here rather than re-reading the clock, so every "now"-relative decision
	// in this one checkout call agrees with the exact instant it was actually
	// let through the open/closed gate.
	results := make([]dto.CheckoutItemResult, 0, len(req.Data))

	var totalbet, totalbayar, totaldiscount, totalkei decimal.Decimal
	var totalpair int

	// limittotal/limitglobal counters live in Redis, not tbl_counter — a key
	// per (period, user/global, tipe, nomor) would be millions of rows under
	// real traffic. The TTL below is a memory cleanup safety net, not a
	// correctness mechanism: nmCounter already embeds idtrxkeluaran, so a new
	// market period always gets a brand new key regardless of TTL.
	limitTTL := 48 * time.Hour
	if closeAt, err := time.ParseInLocation("2006-01-02 15:04:05", pasaran.JadwalTutup, util.LocJakarta); err == nil {
		if remaining := closeAt.Sub(now); remaining > 0 {
			limitTTL = remaining + 24*time.Hour // buffer for late settlement/audit reads after close
		}
	}

	// Redis increments happen outside the Postgres tx below, so they don't
	// get undone for free by tx.Rollback(ctx) if something later in this
	// call fails. Track every accepted increment here and compensate them
	// all in the deferred cleanup unless the tx actually commits.
	type redisLimitIncrement struct {
		key    string
		amount int64
	}
	var redisIncrements []redisLimitIncrement
	committed := false
	defer func() {
		if committed {
			return
		}
		// Detached context: the request ctx may already be canceled/expired
		// on this path (that's often *why* we're unwinding), but the
		// compensating decrement still has to go through.
		cleanupCtx := context.Background()
		for _, inc := range redisIncrements {
			if derr := util.DecrementCounterRedis(cleanupCtx, inc.key, inc.amount); derr != nil {
				connection.Log.Error("failed to roll back redis limit counter after checkout failure",
					zap.String("key", inc.key), zap.Int64("amount", inc.amount), zap.Error(derr))
			}
		}
	}()

	for _, item := range req.Data {
		limittotal, limitglobal := resolveLimits(pasaran, item.Permainan)
		bet := int64(item.Bet)

		nomorForLimit := stripAsterisks(item.Nomor)

		memberKey := "limittotal_" + strings.ToLower(req.Company) + "_" + strconv.Itoa(idtrxkeluaran) +
			"_" + username + "_" + item.Permainan + "_" + item.Nomor
		memberSeed := func(seedCtx context.Context) (int64, error) {
			return detailRepo.SumBet(seedCtx, req.Company, idtrxkeluaran, trxkeluaran.Datekeluaran, username, item.Permainan, nomorForLimit)
		}
		currentMember, okMember, err := util.CheckAndIncrementLimitRedis(ctx, memberKey, bet, limittotal.IntPart(), limitTTL, memberSeed)
		if err != nil {
			return dto.CheckoutResponse{}, err
		}
		if !okMember {
			// currentMember is what's already accumulated on this key *before*
			// this bet — the headroom left is limit minus that, not minus the
			// (too-large) bet that just got rejected.
			sisa := limittotal.Sub(decimal.NewFromInt(currentMember))
			if sisa.IsNegative() {
				sisa = decimal.Zero
			}
			results = append(results, dto.CheckoutItemResult{
				ID: item.ID, Status: "rejected",
				Reason: "melebihi limittotal", Sisalimit: sisa,
			})
			continue
		}
		redisIncrements = append(redisIncrements, redisLimitIncrement{memberKey, bet})

		globalKey := "limitglobal_" + strings.ToLower(req.Company) + "_" + strconv.Itoa(idtrxkeluaran) +
			"_" + item.Permainan + "_" + item.Nomor
		globalSeed := func(seedCtx context.Context) (int64, error) {
			return detailRepo.SumBet(seedCtx, req.Company, idtrxkeluaran, trxkeluaran.Datekeluaran, "", item.Permainan, nomorForLimit)
		}
		currentGlobal, okGlobal, err := util.CheckAndIncrementLimitRedis(ctx, globalKey, bet, limitglobal.IntPart(), limitTTL, globalSeed)
		if err != nil {
			return dto.CheckoutResponse{}, err
		}
		if !okGlobal {
			if err := util.DecrementCounterRedis(ctx, memberKey, bet); err != nil {
				return dto.CheckoutResponse{}, err
			}
			redisIncrements = redisIncrements[:len(redisIncrements)-1]
			sisa := limitglobal.Sub(decimal.NewFromInt(currentGlobal))
			if sisa.IsNegative() {
				sisa = decimal.Zero
			}
			results = append(results, dto.CheckoutItemResult{
				ID: item.ID, Status: "rejected",
				Reason: "melebihi limitglobal", Sisalimit: sisa,
			})
			continue
		}
		redisIncrements = append(redisIncrements, redisLimitIncrement{globalKey, bet})

		// Recomputed server-side from the live pasaran rates — item.Diskon /
		// item.Win / item.Kei / item.Bayar from the client are display-only
		// and never trusted for what actually gets persisted/paid out.
		payout := calculatePayout(pasaran, item.Permainan, item.Nomor, item.Bet, item.Tipetoto)

		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		detail := &domain.Trxkeluarandetail{
			ID:                   now.Format("0601") + raw,
			IDtrxkeluaran:        idtrxkeluaran,
			IDcomp:               idcomp,
			Datekeluarandetail:   sql.NullTime{Valid: true, Time: now},
			Ipaddress:            ipaddress,
			Username:             username,
			Typegame:             item.Permainan,
			Nomortogel:           nomorForLimit,
			Posisitogel:          item.Tipetoto,
			Bet:                  item.Bet,
			Diskon:               payout.DiscRate,
			Win:                  payout.Win,
			Kei:                  payout.KeiRate,
			Browsertogel:         req.Devicemember,
			Devicetogel:          req.Devicemember,
			Statuskeluarandetail: "RUNNING",
			Betround:             betround,
			Playerinvoice:        playerinvoice,
			Created:              username,
			CreatedAt:            sql.NullTime{Valid: true, Time: now},
		}
		if err := detailRepo.Save(ctx, detail, req.Company); err != nil {
			return dto.CheckoutResponse{}, err
		}

		totalbet = totalbet.Add(decimal.NewFromInt(bet))
		totalbayar = totalbayar.Add(payout.Payout)
		totaldiscount = totaldiscount.Add(payout.Disc)
		totalkei = totalkei.Add(payout.Kei)
		totalpair++

		results = append(results, dto.CheckoutItemResult{ID: item.ID, Status: "accepted"})
	}

	if totalpair > 0 {
		// Serializes concurrent checkouts from the SAME player for the SAME
		// period, so the ExistsByUsername check below can't race with
		// itself — two simultaneous requests both seeing "no prior row" and
		// both incrementing total_member. Scoped to this exact
		// (idtrxkeluaran, username) pair and released automatically when the
		// transaction ends (commit or rollback), so it never blocks a
		// DIFFERENT player's checkout for this same period.
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1, hashtext($2))`, idtrxkeluaran, username); err != nil {
			return dto.CheckoutResponse{}, err
		}

		// Checked before this checkout's own member row is saved below — a
		// repeat checkout by a player who already has a row for this period
		// must not count toward total_member again.
		hadPriorInvoice, err := memberRepo.ExistsByUsername(ctx, idcomp, idtrxkeluaran, username)
		if err != nil {
			return dto.CheckoutResponse{}, err
		}

		member := &domain.Trxkeluaranmember{
			ID:            now.Format("0601") + strings.ReplaceAll(uuid.NewString(), "-", ""),
			IDtrxkeluaran: idtrxkeluaran,
			IDcomp:        idcomp,
			Username:      username,
			Totalbet:      totalbet,
			Totalbayar:    totalbayar,
			Totaldiscount: totaldiscount,
			Totalkei:      totalkei,
			Totalpair:     totalpair,
			Betround:      betround,
			Playerinvoice: playerinvoice,
			Status:        "pending",
			Created:       username,
			CreatedAt:     sql.NullTime{Valid: true, Time: now},
		}
		if err := memberRepo.Save(ctx, member, req.Company); err != nil {
			return dto.CheckoutResponse{}, err
		}

		newMember := 0
		if !hadPriorInvoice {
			newMember = 1
		}
		if err := trxkeluaranTxRepo.IncrementTotals(
			ctx, idcomp, idtrxkeluaran, newMember,
			totalbet, decimal.NewFromInt(int64(totalpair)), totalbayar,
		); err != nil {
			return dto.CheckoutResponse{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.CheckoutResponse{}, err
	}
	committed = true

	// Drop the player's cached member summary so their next "Transaksi" fetch
	// reflects this checkout instead of stale data. Detail rows aren't
	// cached (see trxkeluarandetailService.AllByUsername), so there's
	// nothing to invalidate for those.
	go connection.DeleteRedis(trxkeluaranmemberByUserCacheKey(idcomp, req.PasaranIdcomp, username))

	// This checkout just changed total_member/total_bet/total_pairs/
	// total_payout on the period row (see IncrementTotals above) — drop the
	// agen dashboard's cached copies of it (shared Redis, same key format
	// totagen_api's own trxkeluaran.go writes/invalidates under) so agents
	// see the updated totals instead of a stale snapshot.
	go connection.DeleteRedis(RedisTrxkeluaran + ":" + strings.ToLower(idcomp))
	go connection.DeleteRedis(RedisTrxkeluaranDetail + ":" + strings.ToLower(idcomp) + ":" + strconv.Itoa(idtrxkeluaran))

	return dto.CheckoutResponse{
		Playerinvoice: playerinvoice,
		Totalbayar:    totalbayar,
		Items:         results,
	}, nil
}
