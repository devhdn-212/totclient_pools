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
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	RedisCompanywalletKey = "master:companywallet:all:"
)

type companywalletService struct {
	db   *sql.DB
	repo domain.CompanywalletRepository
}

func NewCompanywalletService(db *sql.DB, repo domain.CompanywalletRepository) domain.CompanywalletService {
	return &companywalletService{
		db:   db,
		repo: repo,
	}
}

func (c companywalletService) All(ctx context.Context, idcomp string) ([]dto.CompanywalletData, error) {
	cached, found, err := connection.GetRedis(RedisCompanywalletKey + idcomp)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CompanywalletData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Company Wallet")
			return data, nil
		}
	}

	curr, err := c.repo.FindAll(ctx, idcomp)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var compwalletData []dto.CompanywalletData
	for _, v := range curr {
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

		compwalletData = append(compwalletData, dto.CompanywalletData{
			ID:      v.ID,
			IDcomp:  v.IDcompany,
			IDcurr:  v.IDcurr,
			Amount:  v.Amount,
			Status:  v.Status,
			Created: createdAt,
			Update:  updatedAt,
		})
	}

	go connection.SetRedis(RedisCompanywalletKey+idcomp, compwalletData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Company Wallet")
	return compwalletData, nil
}

func (c companywalletService) Save(ctx context.Context, req dto.CompanywalletSave, client_admin string) error {
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
	txRepo := repository.NewCompanywalletRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID, req.IDcomp, req.IDcurr)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("Duplicate Entry")
		}
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		date := time.Now().Format("060102")
		walletCode := fmt.Sprintf("%s-%s-%s", req.IDcomp, date, raw)
		compwallet := domain.Companywallet{
			ID:        walletCode,
			IDcompany: req.IDcomp,
			IDcurr:    req.IDcurr,
			Status:    req.Status,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &compwallet)
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
			return errors.New("Wallet not found")
		}

		flag.ID = req.ID
		flag.IDcompany = req.IDcomp
		flag.IDcurr = req.IDcurr
		flag.Status = req.Status
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = c.repo.Update(ctx, &flag); err != nil {
			fmt.Println("Error update : ", err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisCompanywalletKey + req.IDcomp)
	return nil
}
