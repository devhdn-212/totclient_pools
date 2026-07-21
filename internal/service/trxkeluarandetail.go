package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/repository"
	"github.com/devhdn-212/totagen_api/internal/util"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

const (
	RedisTrxkeluarandetail = "agen:trxkeluarandetail"
)

type trxkeluarandetailService struct {
	db   *pgxpool.Pool
	repo domain.TrxkeluarandetailRepository
}

func NewTrxkeluarandetailService(db *pgxpool.Pool, repo domain.TrxkeluarandetailRepository) domain.TrxkeluarandetailService {
	return &trxkeluarandetailService{
		db:   db,
		repo: repo,
	}
}

func (u *trxkeluarandetailService) All(ctx context.Context, idcomp string, idtrx int) ([]dto.TrxkeluarandetailData, error) {
	cached, found, err := connection.GetRedis(RedisTrxkeluarandetail + ":" + strings.ToLower(idcomp) + ":" + strconv.Itoa(idtrx))
	if err != nil {
		return nil, err
	}
	var record []dto.TrxkeluarandetailData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Trxkeluarandetail")
			return record, nil
		}
	}

	trxkeluarandetail, err := u.repo.FindAll(ctx, idcomp, idtrx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range trxkeluarandetail {
		var datekeluarandetail, createdAt, updatedAt string

		if v.Datekeluarandetail.Valid {
			datekeluarandetail = v.Datekeluarandetail.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
		}
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
		record = append(record, dto.TrxkeluarandetailData{
			ID:                 v.ID,
			IDtrxkeluaran:      v.IDtrxkeluaran,
			Datekeluarandetail: datekeluarandetail,
			Username:           v.Username,
			Typegame:           v.Typegame,
			Nomortogel:         v.Nomortogel,
			Posisitogel:        v.Posisitogel,
			Bet:                v.Bet,
			Diskon:             v.Diskon,
			Kei:                v.Kei,
			Win:                v.Win,
			Created:            createdAt,
			Update:             updatedAt,
		})
	}
	go connection.SetRedis(RedisTrxkeluarandetail+":"+strings.ToLower(idcomp)+":"+strconv.Itoa(idtrx), record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Trxkeluarandetail")
	return record, nil
}

func (u *trxkeluarandetailService) Save(ctx context.Context, req dto.TrxkeluarandetailSave, client, idcomp string) error {
	// Start Transaction native pgx v5
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewTrxkeluarandetailRepository(txExec)

	flag, err := txRepo.FindByID(ctx, idcomp, req.ID, req.IDtrxkeluaran)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		field_column := "tbl_trx_keluarantogel_detail_" + strings.ToLower(idcomp) + "_" + strconv.Itoa(req.IDtrxkeluaran)
		idcounter, err := util.GetNextCounterManualTx(ctx, tx, field_column)
		if err != nil {
			return err
		}

		playerinvoice := fmt.Sprintf("%d%d", req.IDtrxkeluaran, idcounter)
		playerinvoice_temp, err := strconv.Atoi(playerinvoice)

		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		date := now.Format("0601")
		idtrxkeluarandetail := fmt.Sprintf("%s%s", date, raw)

		flag.ID = idtrxkeluarandetail
		flag.IDtrxkeluaran = req.IDtrxkeluaran
		flag.IDcomp = idcomp
		flag.Datekeluarandetail = sql.NullTime{Valid: true, Time: now}
		flag.Username = "lldsdsb0013794"
		flag.Typegame = idcomp
		flag.Nomortogel = idcomp
		flag.Posisitogel = idcomp
		flag.Bet = 10000
		flag.Diskon = decimal.NewFromFloat(0.05)
		flag.Win = decimal.NewFromFloat(0.05)
		flag.Kei = decimal.NewFromFloat(0.05)
		flag.Playerinvoice = playerinvoice_temp
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

	go connection.DeleteRedis(RedisTrxkeluarandetail + ":" + strings.ToLower(idcomp) + ":" + strconv.Itoa(req.IDtrxkeluaran))
	return nil
}
