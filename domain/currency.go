package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibermaster_api/dto"
)

type Currency struct {
	ID        string       `db:"idcurr"`
	Type      string       `db:"typecurr"`
	Status    string       `db:"status"`
	Created   string       `db:"createcurr"`
	CreatedAt sql.NullTime `db:"createdatecurr"`
	Update    string       `db:"updatecurr"`
	UpdateAt  sql.NullTime `db:"updatedatecurr"`
}

type CurrencyRepository interface {
	FindAll(ctx context.Context) ([]Currency, error)
	FindSelect(ctx context.Context) ([]Currency, error)
	FindByID(ctx context.Context, id string) (Currency, error)
	Save(ctx context.Context, cur *Currency) error
	Update(ctx context.Context, cur *Currency) error
}
type CurrencyService interface {
	All(ctx context.Context) ([]dto.CurrData, error)
	Select(ctx context.Context) ([]dto.CurrSelect, error)
	Save(ctx context.Context, req dto.CurrSave, client string) error
}
