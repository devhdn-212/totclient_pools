package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/repository"
	"gofibergocu/internal/util"
	"time"
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
	exec := goqu.New("default", c.db)
	repo := repository.NewCustomer(exec)

	customers, err := repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var customerData []dto.CustomerData
	for _, v := range customers {
		customerData = append(customerData, dto.CustomerData{
			ID:   v.ID,
			Code: v.Code,
			Name: v.Name,
		})
	}
	return customerData, nil
}
func (c customerService) Create(ctx context.Context, req dto.CreateCustomerRequest) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txExec := goqu.NewTx("default", tx)
	repo := repository.NewCustomer(txExec)

	flag, err := repo.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if flag.ID != "" {
		return errors.New("code already exists")
	}
	customer := domain.Customer{
		ID:        uuid.NewString(),
		Code:      req.Code,
		Name:      req.Name,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	err = repo.Save(ctx, &customer)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return util.ErrDuplicate
			}
		}
	}
	return tx.Commit()

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

	txExec := goqu.NewTx("default", tx)
	repo := repository.NewCustomer(txExec)

	flag, err := repo.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if flag.ID == "" {
		return errors.New("customer id not found")
	}
	flag.Code = req.Code
	flag.Name = req.Name
	flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now()}

	if err = repo.Update(ctx, &flag); err != nil {
		return err
	}
	return tx.Commit()
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

	txExec := goqu.NewTx("default", tx)
	repo := repository.NewCustomer(txExec)

	flag, err := repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if flag.ID == "" {
		return errors.New("customer id not found")
	}

	if err = repo.Delete(ctx, id); err != nil {
		return err
	}
	return tx.Commit()
}
func (c customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	exec := goqu.New("default", c.db)
	repo := repository.NewCustomer(exec)

	flag, err := repo.FindByID(ctx, id)
	if err != nil {
		return dto.CustomerData{}, err
	}
	if flag.ID == "" {
		return dto.CustomerData{}, errors.New("customer id not found")
	}
	return dto.CustomerData{
		ID:   flag.ID,
		Code: flag.Code,
		Name: flag.Name,
	}, nil
}
