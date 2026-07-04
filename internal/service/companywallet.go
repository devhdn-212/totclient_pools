package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	RedisCompanywalletKey = "master:companywallet:all:"
)

type companywalletService struct {
	db   *pgxpool.Pool
	repo domain.CompanywalletRepository
}

func NewCompanywalletService(db *pgxpool.Pool, repo domain.CompanywalletRepository) domain.CompanywalletService {
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

	var compwalletData []dto.CompanywalletData
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
	// Start transaction pgx v5
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer rollback
	defer tx.Rollback(ctx)

	// Executor transaksi native pgx
	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewCompanywalletRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID, req.IDcomp, req.IDcurr)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

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
			CreatedAt: sql.NullTime{Valid: true, Time: now},
		}

		err = txRepo.Save(ctx, &compwallet)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
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
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		// Gunakan txRepo agar berada dalam scope transaksi yang sama
		if err = txRepo.Update(ctx, &flag); err != nil {
			return err
		}
	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisCompanywalletKey + req.IDcomp)
	return nil
}
