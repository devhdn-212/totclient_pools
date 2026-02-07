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
	RedisCustomerAllKey = "customers:all"
	RedisCustomerDetail = "customers:detail:"
)

type customerService struct {
	db   *sql.DB
	repo domain.CustomerRepository
}

func NewCustomerService(db *sql.DB, repo domain.CustomerRepository) domain.CustomerService {
	return &customerService{
		db:   db,
		repo: repo,
	}
}

func (c customerService) Index(ctx context.Context) ([]dto.CustomerData, error) {
	cached, found, err := connection.GetRedis(RedisCustomerAllKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CustomerData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Customer")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	customers, err := c.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var customerData []dto.CustomerData
	for _, v := range customers {
		var createdAt, updatedAt string
		if v.CreatedAt.Valid {
			createdAt = v.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if v.UpdateAt.Valid {
			updatedAt = v.UpdateAt.Time.Format("2006-01-02 15:04:05")
		}
		customerData = append(customerData, dto.CustomerData{
			ID:        v.ID,
			Code:      v.Code,
			Name:      v.Name,
			CreatedAt: createdAt, // langsung sql.NullTime
			UpdateAt:  updatedAt,
		})
	}

	go connection.SetRedis(RedisCustomerAllKey, customerData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Customer")
	return customerData, nil
}
func (c customerService) Create(ctx context.Context, req dto.CreateCustomerRequest) (err error) {
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
	txRepo := repository.NewCustomerRepository(txExec)
	flag, err := txRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if flag.Code != "" {
		return errors.New("code already exists")
	}
	customer := domain.Customer{
		ID:        uuid.NewString(),
		Code:      req.Code,
		Name:      req.Name,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	err = txRepo.Save(ctx, &customer)
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
	go connection.DeleteRedis(RedisCustomerAllKey)
	return nil
}
func (c customerService) Update(ctx context.Context, req dto.UpdateCustomerRequest) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	flag, err := c.repo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if flag.ID == "" {
		return errors.New("customer id not found")
	}
	flag.Name = req.Name
	flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now()}

	if err = c.repo.Update(ctx, &flag); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	go connection.DeleteRedis(RedisCustomerAllKey)
	go connection.DeleteRedis(RedisCustomerDetail + req.ID)
	return nil
}
func (c customerService) Delete(ctx context.Context, id string) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	flag, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if flag.ID == "" {
		return errors.New("customer id not found")
	}

	if err = c.repo.Delete(ctx, id); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	go connection.DeleteRedis(RedisCustomerAllKey)
	go connection.DeleteRedis(RedisCustomerDetail + id)
	return nil
}
func (c customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	redisdetail := RedisCustomerDetail + id

	cached, found, err := connection.GetRedis(redisdetail)
	if err == nil && found {
		var dto dto.CustomerData
		if err := json.Unmarshal([]byte(cached), &dto); err == nil {
			log.Info("Returning data from Redis - Customer/%s", id)
			return dto, nil
		}
	}

	flag, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return dto.CustomerData{}, err
	}
	if flag.ID == "" {
		return dto.CustomerData{}, errors.New("customer id not found")
	}

	result := dto.CustomerData{
		ID:   flag.ID,
		Code: flag.Code,
		Name: flag.Name,
	}

	_ = connection.SetRedis(redisdetail, result, 60*time.Minute)
	log.Info("Returning data Database - Customer/$s", id)
	return result, nil
}
