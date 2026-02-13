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
	RedisBankKey       = "master:bank:all"
	RedisBankSelectKey = "master:bank:select"
)

type bankService struct {
	db   *sql.DB
	repo domain.BankRepository
}

func NewBankService(db *sql.DB, repo domain.BankRepository) domain.BankService {
	return &bankService{
		db:   db,
		repo: repo,
	}
}

func (b bankService) All(ctx context.Context) ([]dto.BankData, error) {
	cached, found, err := connection.GetRedis(RedisBankKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.BankData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Bank")
			return data, nil
		}
	}

	bank, err := b.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var bankData []dto.BankData
	for _, v := range bank {
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

		bankData = append(bankData, dto.BankData{
			ID:      v.ID,
			Type:    v.Type,
			Name:    v.Name,
			Status:  v.Status,
			Created: createdAt,
			Update:  updatedAt,
		})
	}
	go connection.SetRedis(RedisBankKey, bankData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Bank")
	return bankData, nil
}

func (b bankService) Select(ctx context.Context) ([]dto.BankSelect, error) {
	cached, found, err := connection.GetRedis(RedisBankSelectKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.BankSelect
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Bank Select")
			return data, nil
		}
	}

	curr, err := b.repo.FindSelect(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var bankSelect []dto.BankSelect
	for _, v := range curr {
		bankSelect = append(bankSelect, dto.BankSelect{
			ID: v.ID,
		})
	}

	go connection.SetRedis(RedisBankSelectKey, bankSelect, 60*time.Minute)
	connection.Log.Info("Returning data Database - Bank Select")
	return bankSelect, nil
}

func (b bankService) Save(ctx context.Context, req dto.BankSave, client_admin string) error {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewBankRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("Duplicate Entry")
		}

		bank := domain.Bank{
			ID:        req.ID,
			Type:      req.Type_bank,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &bank)
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
			return errors.New("Currency not found")
		}

		flag.ID = req.ID
		flag.Type = req.Type_bank
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = b.repo.Update(ctx, &flag); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisBankKey)
	go connection.DeleteRedis(RedisBankSelectKey)
	return nil
}
