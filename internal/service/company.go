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
	RedisCompanyKey = "master:company:all"
)

type companyService struct {
	db   *sql.DB
	repo domain.CompanyRepository
}

func NewCompanyService(db *sql.DB, repo domain.CompanyRepository) domain.CompanyService {
	return &companyService{
		db:   db,
		repo: repo,
	}
}
func (c companyService) All(ctx context.Context) ([]dto.CompanyData, error) {
	cached, found, err := connection.GetRedis(RedisCompanyKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CompanyData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Company")
			return data, nil
		}
	}

	curr, err := c.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var compData []dto.CompanyData
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

		compData = append(compData, dto.CompanyData{
			ID:        v.ID,
			IDcurrdef: v.IDcurrdef,
			Name:      v.Name,
			Endjoin:   createdAt,
			Status:    v.Status,
			Created:   createdAt,
			Update:    updatedAt,
		})
	}

	go connection.SetRedis(RedisCompanyKey, compData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Company")
	return compData, nil
}

func (c companyService) Save(ctx context.Context, req dto.CompanySave, client_admin string) error {
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
	txRepo := repository.NewCompanyRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("Duplicate Entry")
		}
		comp := domain.Company{
			ID:        req.ID,
			IDcurrdef: req.IDcurr,
			Name:      req.Name,
			Status:    req.Status,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &comp)
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
			return errors.New("Company not found")
		}

		flag.ID = req.ID
		flag.IDcurrdef = req.IDcurr
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

		if err = c.repo.Update(ctx, &flag); err != nil {
			fmt.Println(err)
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go connection.DeleteRedis(RedisCompanyKey)
	return nil
}
