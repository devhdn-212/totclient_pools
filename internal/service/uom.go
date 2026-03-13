package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/dto"
	"github.com/devhdn-212/gofibergoqu_master/internal/connection"
	"github.com/devhdn-212/gofibergoqu_master/internal/repository"
	"github.com/devhdn-212/gofibergoqu_master/internal/util"
	"github.com/gofiber/fiber/v2/log"
	"github.com/lib/pq"
)

const (
	RedisUomAllKey    = "master:uom:all"
	RedisUomSelectKey = "master:uom:select"
)

type uomService struct {
	db   *sql.DB
	repo domain.UomRepository
}

func NewUomService(db *sql.DB, repo domain.UomRepository) domain.UomService {
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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	for _, v := range curr {
		record = append(record, dto.UomData{
			ID:      v.ID,
			Name:    v.Name,
			Status:  v.Status,
			Created: util.NtToStr(v.CreatedAt, v.Created, loc),
			Update:  util.NtToStr(v.UpdateAt, v.Update, loc),
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
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewUomRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("Duplicate Entry")
		}

		curr := domain.Uom{
			ID:        req.ID,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &curr)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
		if err = tx.Commit(); err != nil {
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
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = u.repo.Update(ctx, &flag); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisUomAllKey)
	go connection.DeleteRedis(RedisUomSelectKey)
	return nil
}
