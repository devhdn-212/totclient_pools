package domain

import (
	"context"
	"database/sql"
	"gofibergocu/dto"

	"github.com/shopspring/decimal"
)

type Company struct {
	ID        string          `db:"idcompany"`
	IDcurrdef string          `db:"idcurrdef"`
	Name      string          `db:"compname"`
	Endjoin   sql.NullTime    `db:"endjoin"`
	Amount    decimal.Decimal `db:"amountcomp"`
	Status    string          `db:"compstatus"`
	Created   string          `db:"createcomp"`
	CreatedAt sql.NullTime    `db:"createdatecomp"`
	Update    string          `db:"updatecomp"`
	UpdateAt  sql.NullTime    `db:"updatedatecomp"`
}

type CompanyRepository interface {
	FindAll(ctx context.Context) ([]Company, error)
	FindByID(ctx context.Context, id string) (Company, error)
	Save(ctx context.Context, cur *Company) error
	Update(ctx context.Context, cur *Company) error
}
type CompanyService interface {
	All(ctx context.Context) ([]dto.CompanyData, error)
	Save(ctx context.Context, req dto.CompanySave, client string) error
}
