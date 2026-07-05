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
	RedisGroupcompanyAllKey    = "master:groupcompany:all"
	RedisGroupcompanySelectKey = "master:groupcompany:select"
)

type groupCompanyService struct {
	db   *pgxpool.Pool
	repo domain.GroupcompanyRepository
}

func NewGroupCompanyService(db *pgxpool.Pool, repo domain.GroupcompanyRepository) domain.GroupcompanyService {
	return &groupCompanyService{
		db:   db,
		repo: repo,
	}
}

func (g *groupCompanyService) All(ctx context.Context) ([]dto.GroupcompanyData, error) {
	cached, found, err := connection.GetRedis(RedisGroupcompanyAllKey)
	if err != nil {
		return nil, err
	}
	var record []dto.GroupcompanyData
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Groupcompany")
			return record, nil
		}
	}

	curr, err := g.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

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
		record = append(record, dto.GroupcompanyData{
			ID:      v.ID,
			Name:    v.Name,
			Status:  v.Status,
			Created: createdAt,
			Update:  updatedAt,
		})
	}
	go connection.SetRedis(RedisGroupcompanyAllKey, record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Groupcompany")
	return record, nil
}
func (g *groupCompanyService) Select(ctx context.Context) ([]dto.GroupcompanySelect, error) {
	cached, found, err := connection.GetRedis(RedisGroupcompanySelectKey)
	if err != nil {
		return nil, err
	}
	var record []dto.GroupcompanySelect
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			connection.Log.Info("Returning data from Redis - Groupcompany Select")
			return record, nil
		}
	}

	curr, err := g.repo.FindSelect(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, v := range curr {
		record = append(record, dto.GroupcompanySelect{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	go connection.SetRedis(RedisGroupcompanySelectKey, record, 24*time.Hour)
	connection.Log.Info("Returning data Database - Groupcompany")
	return record, nil
}

func (g *groupCompanyService) Save(ctx context.Context, req dto.GroupcompanySave, client string) error {
	// Start Transaction native pgx v5
	tx, err := g.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback jika terjadi panic atau error sebelum commit
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewGroupcompanyRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}
	fmt.Println("Flag Groupcompany: ", req.ID, flag.ID)
	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != "" {
			return util.ErrDuplicate
		}

		groupcompany := domain.Groupcompany{
			ID:        req.ID,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client,
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}
		err = txRepo.Save(ctx, &groupcompany)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
	} else {
		if flag.ID == "" {
			return errors.New("Groupcompany not found")
		}

		flag.ID = req.ID
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		// Gunakan txRepo agar tetap dalam satu scope transaksi
		if err = txRepo.Update(ctx, &flag); err != nil {
			fmt.Println("Error update Groupcompany: ", err)
			return err
		}
	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisGroupcompanyAllKey)
	go connection.DeleteRedis(RedisGroupcompanySelectKey)
	return nil
}
