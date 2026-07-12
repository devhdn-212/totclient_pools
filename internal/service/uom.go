package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/repository"
	"github.com/devhdn-212/totagen_api/internal/util"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisUomAllKey    = "master:uom:all"
	RedisUomSelectKey = "master:uom:select"
)

type uomService struct {
	db   *pgxpool.Pool
	repo domain.UomRepository
}

func NewUomService(db *pgxpool.Pool, repo domain.UomRepository) domain.UomService {
	return &uomService{
		db:   db,
		repo: repo,
	}
}

// All implements [domain.UomService].
func (u *uomService) All(ctx context.Context) ([]dto.UomData, error) {
	cached, found, err := connection.GetRedis(RedisUomAllKey)
	if err != nil {
		return nil, err
	}
	var record []dto.UomData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - UOM")
			return record, nil
		}
	}

	curr, err := u.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range curr {
		record = append(record, dto.UomData{
			ID:      v.ID,
			Name:    v.Name,
			Status:  v.Status,
			Created: util.NtToStr(v.CreatedAt, v.Created, util.LocJakarta),
			Update:  util.NtToStr(v.UpdateAt, v.Update, util.LocJakarta),
		})
	}
	go connection.SetRedis(RedisUomAllKey, record, 24*time.Hour)
	connection.Log.Info("Returning data Database - UOM")
	return record, nil
}
func (u *uomService) Select(ctx context.Context) ([]dto.UomSelect, error) {
	cached, found, err := connection.GetRedis(RedisUomSelectKey)
	if err != nil {
		return nil, err
	}
	var record []dto.UomSelect
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - UOM")
			return record, nil
		}
	}

	curr, err := u.repo.FindSelect(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, v := range curr {
		record = append(record, dto.UomSelect{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	go connection.SetRedis(RedisUomSelectKey, record, 24*time.Hour)
	connection.Log.Info("Returning data Database - UOM")
	return record, nil
}

// Save implements [domain.UomService].
func (u *uomService) Save(ctx context.Context, req dto.UomSave, client string) error {
	// Start Transaction native pgx v5
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewUomRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != "" {
			return util.ErrDuplicate
		}

		uom := domain.Uom{
			ID:        req.ID,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client,
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}
		err = txRepo.Save(ctx, &uom)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
	} else {
		if flag.ID == "" {
			return errors.New("UOM not found")
		}

		flag.ID = req.ID
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		// Gunakan txRepo agar tetap dalam satu scope transaksi
		if err = txRepo.Update(ctx, &flag); err != nil {
			fmt.Println("Error update UOM: ", err)
			return err
		}
	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisUomAllKey)
	go connection.DeleteRedis(RedisUomSelectKey)
	return nil
}
