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
	RedisClientruleAllKey    = "master:clientrule:all"
	RedisClientruleSelectKey = "master:clientrule:select"
)

type clientruleService struct {
	db   *sql.DB
	repo domain.ClientruleRepository
}

func NewClientruleService(db *sql.DB, repo domain.ClientruleRepository) domain.ClientruleService {
	return &clientruleService{
		db:   db,
		repo: repo,
	}
}
func (c clientruleService) All(ctx context.Context) ([]dto.ClientruleData, error) {
	cached, found, err := connection.GetRedis(RedisClientruleAllKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.ClientruleData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Clientrule")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	clr, err := c.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var clientruleData []dto.ClientruleData
	for _, v := range clr {
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

		clientruleData = append(clientruleData, dto.ClientruleData{
			ID:      v.ID,
			Name:    v.Name,
			Rule:    v.Rule,
			Created: createdAt,
			Update:  updatedAt,
		})
	}

	go connection.SetRedis(RedisClientruleAllKey, clientruleData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Clientrule")
	return clientruleData, nil
}

func (c clientruleService) Select(ctx context.Context) ([]dto.ClientruleSelect, error) {
	cached, found, err := connection.GetRedis(RedisClientruleAllKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.ClientruleSelect
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Clientrule Select")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	admins, err := c.repo.FindSelect(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var clientruleSelect []dto.ClientruleSelect
	for _, v := range admins {
		clientruleSelect = append(clientruleSelect, dto.ClientruleSelect{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	go connection.SetRedis(RedisClientruleAllKey, clientruleSelect, 60*time.Minute)
	connection.Log.Info("Returning data Database - Clientrule Select")
	return clientruleSelect, nil
}

func (c clientruleService) Save(ctx context.Context, req dto.ClientruleSave, client string) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewClientruleRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("ID already exists")
		}
		clientrule := domain.Clientrule{
			ID:        req.ID,
			Name:      req.Name,
			Rule:      req.Rule,
			Created:   client,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &clientrule)
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
		flag.Update = client
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = c.repo.Update(ctx, &flag); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisClientruleAllKey)
	go connection.DeleteRedis(RedisClientruleSelectKey)
	return nil
}
