package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totmaster_api/dto"

	"github.com/shopspring/decimal"
)

type Companywallet struct {
	ID        string          `db:"idcompwallet"`
	IDcompany string          `db:"idcompany"`
	IDcurr    string          `db:"idcurr"`
	Amount    decimal.Decimal `db:"amountcompwallet"`
	Status    string          `db:"compwalletstatus"`
	Created   string          `db:"createcompwallet"`
	CreatedAt sql.NullTime    `db:"createdatecompwallet"`
	Update    string          `db:"updatecompwallet"`
	UpdateAt  sql.NullTime    `db:"updatedatecompwallet"`
}

type CompanywalletRepository interface {
	FindAll(ctx context.Context, id string) ([]Companywallet, error)
	FindByID(ctx context.Context, id, idcomp, idcurr string) (Companywallet, error)
	Save(ctx context.Context, cur *Companywallet) error
	Update(ctx context.Context, cur *Companywallet) error
}

type CompanywalletService interface {
	All(ctx context.Context, idcomp string) ([]dto.CompanywalletData, error)
	Save(ctx context.Context, req dto.CompanywalletSave, client string) error
}
