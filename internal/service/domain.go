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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisDomainAllKey = "master:domain:all"
)

type domainService struct {
	db   *pgxpool.Pool
	repo domain.DomainRepository
}

func NewDomainService(db *pgxpool.Pool, repo domain.DomainRepository) domain.DomainService {
	return &domainService{
		db:   db,
		repo: repo,
	}
}
func (d domainService) All(ctx context.Context) ([]dto.DomainData, error) {
	cached, found, err := connection.GetRedis(RedisDomainAllKey)
	if err != nil {
		return nil, err
	}
	var data []dto.DomainData
	if found {
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Domain")
			return data, nil
		}
	}

	dm, err := d.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range dm {
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

		data = append(data, dto.DomainData{
			ID:      v.ID,
			Type:    v.Type,
			Name:    v.Name,
			Status:  v.Status,
			Created: createdAt,
			Update:  updatedAt,
		})
	}
	go connection.SetRedis(RedisDomainAllKey, data, 60*time.Minute)
	connection.Log.Info("Returning data Database - Domain")
	return data, nil
}

func (d domainService) Save(ctx context.Context, req dto.DomainSave, client string) error {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewDomainRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != "" {
			return util.ErrDuplicate
		}

		dm := domain.Domain{
			ID:        uuid.NewString(),
			Type:      req.Typedomain,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client,
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}
		err = txRepo.Save(ctx, &dm)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
	} else {
		if flag.ID == "" {
			return fmt.Errorf("Domain conf toto %w", util.ErrNotFound)
		}

		flag.Type = req.Typedomain
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		// Gunakan txRepo agar perubahan terikat pada transaksi
		if err = txRepo.Update(ctx, &flag); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisDomainAllKey)
	return nil
}
