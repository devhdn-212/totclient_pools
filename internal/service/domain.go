package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/repository"
	"gofibergocu/internal/util"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	RedisDomainAllKey = "master:domain:all"
)

type domainService struct {
	db   *sql.DB
	repo domain.DomainRepository
}

func NewDomainService(db *sql.DB, repo domain.DomainRepository) domain.DomainService {
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

	if found {
		var data []dto.DomainData
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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var domainData []dto.DomainData
	for _, v := range dm {
		var createdAt, updatedAt string
		if v.CreatedAt.Valid {
			createdAt = v.Created + ", " + v.CreatedAt.Time.In(loc).Format("2006-01-02 15:04:05")
		}
		if v.UpdateAt.Valid {
			if v.Update != "" {
				updatedAt = v.Update + ", " + v.UpdateAt.Time.In(loc).Format("2006-01-02 15:04:05")
			} else {
				updatedAt = ""
			}
		}

		domainData = append(domainData, dto.DomainData{
			ID:      v.ID,
			Type:    v.Type,
			Name:    v.Name,
			Status:  v.Status,
			Created: createdAt,
			Update:  updatedAt,
		})
	}
	go connection.SetRedis(RedisDomainAllKey, domainData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Domain")
	return domainData, nil
}

func (d domainService) Save(ctx context.Context, req dto.DomainSave, client string) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewDomainRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("Duplicate Entry")
		}

		dm := domain.Domain{
			ID:        uuid.NewString(),
			Type:      req.Typedomain,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &dm)
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
			return errors.New("Domain not found")
		}

		flag.ID = req.ID
		flag.Type = req.Typedomain
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = d.repo.Update(ctx, &flag); err != nil {
			log.Error(err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisDomainAllKey)
	return nil
}
