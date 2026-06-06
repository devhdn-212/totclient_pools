package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibermaster_api/dto"
)

type Bank struct {
	ID        string         `db:"idbank"`
	Type      string         `db:"typebank"`
	Name      string         `db:"nmbank"`
	Status    string         `db:"bankstatus"`
	Created   string         `db:"createbank"`
	CreatedAt sql.NullTime   `db:"createdatebank"`
	Update    sql.NullString `db:"updatebank"`
	UpdateAt  sql.NullTime   `db:"updatedatebank"`
}

type BankRepository interface {
	FindAll(ctx context.Context) ([]Bank, error)
	FindSelect(ctx context.Context) ([]Bank, error)
	FindByID(ctx context.Context, id string) (Bank, error)
	Save(ctx context.Context, cur *Bank) error
	Update(ctx context.Context, cur *Bank) error
}
type BankService interface {
	All(ctx context.Context) ([]dto.BankData, error)
	Select(ctx context.Context) ([]dto.BankSelect, error)
	Save(ctx context.Context, req dto.BankSave, client string) error
}
