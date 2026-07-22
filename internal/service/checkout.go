package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/repository"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

	now := util.GetNowJakarta()
	results := make([]dto.CheckoutItemResult, 0, len(req.Data))

	var totalbet, totalbayar, totaldiscount, totalkei decimal.Decimal
	var totalpair int

	for _, item := range req.Data {
		limittotal, limitglobal := resolveLimits(pasaran, item.Permainan)
		bet := int64(item.Bet)

		memberKey := "limittotal_" + strings.ToLower(req.Company) + "_" + strconv.Itoa(idtrxkeluaran) +
			"_" + username + "_" + item.Permainan + "_" + item.Nomor
		currentMember, okMember, err := util.CheckAndIncrementLimitTx(ctx, tx, memberKey, bet, limittotal.IntPart())
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

		globalKey := "limitglobal_" + strings.ToLower(req.Company) + "_" + strconv.Itoa(idtrxkeluaran) +
			"_" + item.Permainan + "_" + item.Nomor
		currentGlobal, okGlobal, err := util.CheckAndIncrementLimitTx(ctx, tx, globalKey, bet, limitglobal.IntPart())
		if err != nil {
			return dto.CheckoutResponse{}, err
		}
		if !okGlobal {
			if err := util.DecrementCounterTx(ctx, tx, memberKey, bet); err != nil {
				return dto.CheckoutResponse{}, err
			}
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
			Nomortogel:           stripAsterisks(item.Nomor),
			Posisitogel:          item.Tipetoto,
			Bet:                  item.Bet,
			Diskon:               payout.Disc,
			Win:                  payout.Win,
			Kei:                  payout.Kei,
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
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.CheckoutResponse{}, err
	}

	// Drop the player's cached "Transaksi" list so their next fetch reflects
	// this checkout instead of stale data.
	go connection.DeleteRedis(trxkeluarandetailByUserCacheKey(idcomp, idtrxkeluaran, username))
	go connection.DeleteRedis(trxkeluaranmemberByUserCacheKey(idcomp, idtrxkeluaran, username))

	return dto.CheckoutResponse{
		Playerinvoice: playerinvoice,
		Totalbayar:    totalbayar,
		Items:         results,
	}, nil
}
