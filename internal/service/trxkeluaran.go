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
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisTrxkeluaran = "agen:trxkeluaran"
)

type trxkeluaranService struct {
	db   *pgxpool.Pool
	repo domain.TrxkeluaranRepository
}

func NewTrxkeluaranService(db *pgxpool.Pool, repo domain.TrxkeluaranRepository) domain.TrxkeluaranService {
	return &trxkeluaranService{
		db:   db,
		repo: repo,
	}
}

func (u *trxkeluaranService) All(ctx context.Context, idcomp string) ([]dto.TrxkeluaranData, error) {
	cached, found, err := connection.GetRedis(RedisTrxkeluaran + ":" + strings.ToLower(idcomp))
	if err != nil {
		return nil, err
	}
	var record []dto.TrxkeluaranData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Trxkeluaran")
			return record, nil
		}
	}

	trxkeluaran, err := u.repo.FindAllRunning(ctx, idcomp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range trxkeluaran {

		record = append(record, dto.TrxkeluaranData{
			ID:            v.ID,
			IDcompasaran:  v.IDcomppasaran,
			Nmpasaran:     v.Nmpasaran,
			Datekeluaran:  v.Datekeluaran.Format("2006-01-02"),
			Periode:       v.Keluaranperiode,
			Total_member:  v.Total_member,
			Total_bet:     v.Total_bet,
			Total_buangan: v.Total_buangan,
		})
	}
	go connection.SetRedis(RedisTrxkeluaran+":"+strings.ToLower(idcomp), record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Trxkeluaran")
	return record, nil
}

func (u *trxkeluaranService) Save(ctx context.Context, req dto.TrxkeluaranSave, client, idcomp string) error {
	// Start Transaction native pgx v5
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewTrxkeluaranRepository(txExec)

	flag, err := txRepo.FindByIDByNomorKeluaran(ctx, idcomp, req.IDcompasaran)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != 0 {
			return util.ErrDuplicate
		}

		//periode
		yearperiode := now.Format("06")
		field_column_periode := req.IDcompasaran + "_" + strings.ToLower(idcomp) + "_" + yearperiode
		idperiode, err := util.GetNextCounterManualTx(ctx, tx, field_column_periode)
		if err != nil {
			return err
		}
		idperiode_temp := int(idperiode)

		yearmonth := now.Format("0601")
		field_column := "tbl_trx_keluarantogel" + strings.ToLower(idcomp) + "_" + yearmonth
		idcounter, err := util.GetNextCounterManualTx(ctx, tx, field_column)
		if err != nil {
			return err
		}
		datejam := now.Format("06010215")
		idtrxkeluaran, _ := strconv.Atoi(fmt.Sprintf("%s%d", datejam, idcounter))
		flag.ID = idtrxkeluaran
		flag.IDcomppasaran = req.IDcompasaran
		flag.IDcomp = idcomp
		flag.Yearmonth = now.Format("06-01")
		flag.Keluaranperiode = idperiode_temp
		flag.Datekeluaran = now
		flag.Created = client
		flag.CreatedAt = sql.NullTime{Valid: true, Time: now}

		if err = txRepo.Save(ctx, &flag, idcomp); err != nil {
			fmt.Println("Error New Pasaran: ", err)
			return err
		}
	} else if req.Type == "Edit" {

	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisTrxkeluaran + ":" + strings.ToLower(idcomp))
	return nil
}
