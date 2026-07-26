package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisTrxkeluaran = "agen:trxkeluaran"
	// RedisTrxkeluaranDetail mirrors the agen dashboard's own per-period
	// cache key (see totagen_api's trxkeluaran.go) — keyed by idtrxkeluaran,
	// not just idcomp. checkout.go invalidates this alongside RedisTrxkeluaran
	// since both go stale on the agen side the moment
	// total_member/total_bet/total_pairs/total_payout change here.
	RedisTrxkeluaranDetail = RedisTrxkeluaran + ":detail"
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
