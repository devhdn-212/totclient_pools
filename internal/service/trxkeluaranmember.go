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
	"github.com/devhdn-212/totclient_pools/dto"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/devhdn-212/totclient_pools/internal/repository"
	"github.com/devhdn-212/totclient_pools/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	RedisTrxkeluaranmember       = "agen:trxkeluaranmember"
	RedisTrxkeluaranmemberByUser = "client:trxkeluaranmember"
)

// trxkeluaranmemberByUserCacheKey is keyed by idcomppasaran (not idtrx) so it
// stays valid across period rollovers — checkout.go deletes this exact key
// after a successful checkout.
//
// The ":vN" segment must be bumped every time dto.TrxkeluaranmemberPeriodData's
// shape changes (most recently: v4 switches to aggregated per-period rows
// with totalpair/totalbayar/totalwin). Without it, a pre-existing 24h-TTL
// cache entry unmarshals fine into the new struct and just silently leaves
// the new fields blank — bumping this makes old entries simply miss instead
// of serving stale data.
func trxkeluaranmemberByUserCacheKey(idcomp, idcomppasaran, username string) string {
	return RedisTrxkeluaranmemberByUser + ":v4:" + strings.ToLower(idcomp) + ":" + idcomppasaran + ":" + username
}

type trxkeluaranmemberService struct {
	db   *pgxpool.Pool
	repo domain.TrxkeluaranmemberRepository
}

func NewTrxkeluaranmemberService(db *pgxpool.Pool, repo domain.TrxkeluaranmemberRepository) domain.TrxkeluaranmemberService {
	return &trxkeluaranmemberService{
		db:   db,
		repo: repo,
	}
}

func (u *trxkeluaranmemberService) All(ctx context.Context, idcomp string, idtrx int) ([]dto.TrxkeluaranmemberData, error) {
	cached, found, err := connection.GetRedis(RedisTrxkeluaranmember + ":" + strings.ToLower(idcomp) + ":" + strconv.Itoa(idtrx))
	if err != nil {
		return nil, err
	}
	var record []dto.TrxkeluaranmemberData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Trxkeluaranmember")
			return record, nil
		}
	}

	trxkeluarandetail, err := u.repo.FindAll(ctx, idcomp, idtrx)
	if err != nil {
		connection.Log.Error("trxkeluaranmemberRepository.FindAll failed", zap.Error(err))
		return nil, err
	}

	for _, v := range trxkeluarandetail {
		var createdAt, updatedAt string

		if v.CreatedAt.Valid {
			createdAt = v.Created + ", " + v.CreatedAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
		}
		if v.UpdateAt.Valid {
			if v.Update != "" {
				updatedAt = v.Update + ", " + v.UpdateAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
			} else {
				updatedAt = ""
			}
		}
		record = append(record, dto.TrxkeluaranmemberData{
			ID:            v.ID,
			IDtrxkeluaran: v.IDtrxkeluaran,
			Username:      v.Username,
			Totalbet:      v.Totalbet,
			Totalbayar:    v.Totalbayar,
			Totaldiscount: v.Totaldiscount,
			Totalkei:      v.Totalkei,
			Totalwin:      v.Totalwin,
			Totalpair:     v.Totalpair,
			Betround:      v.Betround,
			Playerinvoice: v.Playerinvoice,
			Status:        v.Status,
			Created:       createdAt,
			Update:        updatedAt,
		})
	}
	go connection.SetRedis(RedisTrxkeluaranmember+":"+strings.ToLower(idcomp)+":"+strconv.Itoa(idtrx), record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Trxkeluaranmember")
	return record, nil
}

func (u *trxkeluaranmemberService) Periods(ctx context.Context, idcomp, idcomppasaran, username string) ([]dto.TrxkeluaranmemberPeriodData, error) {
	cacheKey := trxkeluaranmemberByUserCacheKey(idcomp, idcomppasaran, username)
	cached, found, err := connection.GetRedis(cacheKey)
	if err != nil {
		return nil, err
	}
	var record []dto.TrxkeluaranmemberPeriodData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Trxkeluaranmember (periods)")
			return record, nil
		}
	}

	rows, err := u.repo.FindPeriodsByUsername(ctx, idcomp, idcomppasaran, username)
	if err != nil {
		connection.Log.Error("trxkeluaranmemberRepository.FindPeriodsByUsername failed", zap.Error(err))
		return nil, err
	}

	for _, v := range rows {
		status := "pending"
		if v.Keluarantogel != "" {
			status = "complete"
		}
		record = append(record, dto.TrxkeluaranmemberPeriodData{
			IDtrxkeluaran:    v.IDtrxkeluaran,
			Datekeluaran:     v.Datekeluaran.In(util.LocJakarta).Format("2006-01-02"),
			Keluaranperiode:  v.Keluaranperiode,
			Aliascomppasaran: v.Aliascomppasaran,
			Codecomppasaran:  v.Codecomppasaran,
			Totalpair:        v.Totalpair,
			Totalbayar:       v.Totalbayar,
			Totalwin:         v.Totalwin,
			Status:           status,
		})
	}
	go connection.SetRedis(cacheKey, record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Trxkeluaranmember (periods)")
	return record, nil
}

func (u *trxkeluaranmemberService) Save(ctx context.Context, req dto.TrxkeluaranmemberSave, client, idcomp string) error {
	// Start Transaction native pgx v5
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewTrxkeluaranmemberRepository(txExec)

	flag, err := txRepo.FindByID(ctx, idcomp, req.ID, req.IDtrxkeluaran)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		date := now.Format("06010215")
		idtrxkeluaranmember := fmt.Sprintf("%s%s", date, raw)

		flag.ID = idtrxkeluaranmember
		flag.IDtrxkeluaran = req.IDtrxkeluaran
		flag.IDcomp = idcomp
		flag.Username = "lldsdsb0013794"
		flag.Totalbet = decimal.NewFromInt(1000)
		flag.Totalbayar = decimal.NewFromInt(2000)
		flag.Totaldiscount = decimal.NewFromInt(2000)
		flag.Totalkei = decimal.NewFromInt(2000)
		flag.Totalwin = decimal.NewFromInt(2000)
		flag.Betround = 1
		flag.Playerinvoice = req.Playerinvoice
		flag.Created = client
		flag.CreatedAt = sql.NullTime{Valid: true, Time: now}

		if err = txRepo.Save(ctx, &flag, idcomp); err != nil {
			fmt.Println("Error New Trxkeluarandetail: ", err)
			return err
		}
	} else if req.Type == "Edit" {

	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisTrxkeluaranmember + ":" + strings.ToLower(idcomp) + ":" + strconv.Itoa(req.IDtrxkeluaran))
	return nil
}
