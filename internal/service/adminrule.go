package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/repository"
	"gofibergocu/internal/util"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/lib/pq"
)

const (
	RedisAdminruleAllKey    = "adminrule:all"
	RedisAdminruleSelectKey = "adminrule:select"
)

type adminruleService struct {
	db   *sql.DB
	repo domain.AdminruleRepository
}

func NewAdminruleService(db *sql.DB, repo domain.AdminruleRepository) domain.AdminruleService {
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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var adminruleData []dto.AdminruleData
	for _, v := range admins {
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

	if found {
		var data []dto.AdminruleSelect
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
	var adminruleSelect []dto.AdminruleSelect
	for _, v := range admins {
		adminruleSelect = append(adminruleSelect, dto.AdminruleSelect{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	go connection.SetRedis(RedisAdminruleSelectKey, adminruleSelect, 60*time.Minute)
	connection.Log.Info("Returning data Database - Adminrule Select")
	return adminruleSelect, nil
}
func (a adminruleService) Save(ctx context.Context, req dto.AdminruleSave, client_admin string) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewAdminruleRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("ID already exists")
		}
		adminrule := domain.Adminrule{
			ID:        req.ID,
			Name:      req.Name,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &adminrule)
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
			return errors.New("ID not found")
		}
		flag.Name = req.Name
		flag.Rule = req.Rule
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = a.repo.Update(ctx, &flag); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisAdminruleAllKey)
	go connection.DeleteRedis(RedisAdminruleSelectKey)
	return nil
}
