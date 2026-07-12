package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	RedisAdminruleAllKey    = "master:adminrule:all"
	RedisAdminruleSelectKey = "master:adminrule:select"
)

type adminruleService struct {
	db   *pgxpool.Pool
	repo domain.AdminruleRepository
}

func NewAdminruleService(db *pgxpool.Pool, repo domain.AdminruleRepository) domain.AdminruleService {
	return &adminruleService{
		db:   db,
		repo: repo,
	}
}

func (a adminruleService) All(ctx context.Context) ([]dto.AdminruleData, error) {
	cached, found, err := connection.GetRedis(RedisAdminruleAllKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.AdminruleData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Adminrule")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	admins, err := a.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var adminruleData []dto.AdminruleData
	for _, v := range admins {
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

		adminruleData = append(adminruleData, dto.AdminruleData{
			ID:      v.ID,
			Name:    v.Name,
			Rule:    v.Rule,
			Created: createdAt,
			Update:  updatedAt,
		})
	}

	go connection.SetRedis(RedisAdminruleAllKey, adminruleData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Adminrule")
	return adminruleData, nil
}
func (a adminruleService) Select(ctx context.Context) ([]dto.AdminruleSelect, error) {
	cached, found, err := connection.GetRedis(RedisAdminruleSelectKey)
	if err != nil {
		return nil, err
	}
	var data []dto.AdminruleSelect
	if found {
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Adminrule Select")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	admins, err := a.repo.FindSelect(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range admins {
		data = append(data, dto.AdminruleSelect{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	go connection.SetRedis(RedisAdminruleSelectKey, data, 60*time.Minute)
	connection.Log.Info("Returning data Database - Adminrule Select")
	return data, nil
}
func (a adminruleService) Save(ctx context.Context, req dto.AdminruleSave, client_admin string) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewAdminruleRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// Pakai helper GetNowJakarta()
	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != "" {
			return util.ErrDuplicate
		}
		adminrule := domain.Adminrule{
			ID:        req.ID,
			Name:      req.Name,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}
		err = txRepo.Save(ctx, &adminrule)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
	} else {
		if flag.ID == "" {
			return errors.New("ID not found")
		}
		flag.Name = req.Name
		flag.Rule = req.Rule
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		// Pastikan pakai repository transaksi supaya perubahan masuk dalam scope Tx
		if err = txRepo.Update(ctx, &flag); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisAdminruleAllKey)
	go connection.DeleteRedis(RedisAdminruleSelectKey)
	return nil
}
