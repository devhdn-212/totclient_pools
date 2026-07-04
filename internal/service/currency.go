package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/dto"
	"github.com/devhdn-212/totmaster_api/internal/connection"
	"github.com/devhdn-212/totmaster_api/internal/repository"
	"github.com/devhdn-212/totmaster_api/internal/util"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisCurrAllKey    = "master:curr:all"
	RedisCurrSelectKey = "master:curr:select"
)

type currService struct {
	db   *pgxpool.Pool
	repo domain.CurrencyRepository
}

func NewCurrService(db *pgxpool.Pool, repo domain.CurrencyRepository) domain.CurrencyService {
	return &currService{
		db:   db,
		repo: repo,
	}
}
func (c currService) All(ctx context.Context) ([]dto.CurrData, error) {
	cached, found, err := connection.GetRedis(RedisCurrAllKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CurrData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Currency")
			return data, nil
		}
	}

	curr, err := c.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var currData []dto.CurrData
	for _, v := range curr {
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

		currData = append(currData, dto.CurrData{
			ID:      v.ID,
			Type:    v.Type,
			Status:  v.Status,
			Created: createdAt,
			Update:  updatedAt,
		})
	}
	go connection.SetRedis(RedisCurrAllKey, currData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Currency")
	return currData, nil
}
func (c currService) Select(ctx context.Context) ([]dto.CurrSelect, error) {
	cached, found, err := connection.GetRedis(RedisCurrSelectKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CurrSelect
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Currency Select")
			return data, nil
		}
	}

	curr, err := c.repo.FindSelect(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var currSelect []dto.CurrSelect
	for _, v := range curr {
		currSelect = append(currSelect, dto.CurrSelect{
			ID: v.ID,
		})
	}

	go connection.SetRedis(RedisCurrSelectKey, currSelect, 60*time.Minute)
	connection.Log.Info("Returning data Database - Currency Select")
	return currSelect, nil
}
func (c currService) Save(ctx context.Context, req dto.CurrSave, client_admin string) error {
	// Start Transaction native pgx
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewCurrRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != "" {
			return util.ErrDuplicate
		}

		curr := domain.Currency{
			ID:        req.ID,
			Type:      req.Type_curr,
			Status:    req.Status,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}
		err = txRepo.Save(ctx, &curr)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
	} else {
		if flag.ID == "" {
			return errors.New("Currency not found")
		}

		flag.Type = req.Type_curr
		flag.Status = req.Status
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		// Perbaikan: gunakan txRepo agar tetap dalam transaksi yang sama
		if err = txRepo.Update(ctx, &flag); err != nil {
			fmt.Println(err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisCurrAllKey)
	go connection.DeleteRedis(RedisCurrSelectKey)
	return nil
}
