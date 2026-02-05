package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/util"
	"time"
)

type customerService struct {
	customerRepository domain.CustomerRepository
}

func NewCustomer(customerRepository domain.CustomerRepository) domain.CustomerService {
	return &customerService{
		customerRepository: customerRepository,
	}
}

func (c customerService) Index(ctx context.Context) ([]dto.CustomerData, error) {
	customers, err := c.customerRepository.FindAll(ctx)
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
	flag, err := c.customerRepository.FindByCode(ctx, req.Code)
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
	err = c.customerRepository.Save(ctx, &customer)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return util.ErrDuplicate
			}
		}
	}
	return err

}
func (c customerService) Update(ctx context.Context, req dto.UpdateCustomerRequest) error {
	flag, err := c.customerRepository.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if flag.ID == "" {
		return errors.New("customer id not found")
	}
	flag.Code = req.Code
	flag.Name = req.Name
	flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now()}

	return c.customerRepository.Update(ctx, &flag)
}
func (c customerService) Delete(ctx context.Context, id string) error {
	flag, err := c.customerRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if flag.ID == "" {
		return errors.New("customer id not found")
	}

	return c.customerRepository.Delete(ctx, id)
}
func (c customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	flag, err := c.customerRepository.FindByID(ctx, id)
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
