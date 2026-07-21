package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totclient_api/dto"
)

type Domain struct {
	ID        string       `db:"iddomain"`
	Name      string       `db:"nmdomain"`
	Type      string       `db:"tipedomain"`
	Status    string       `db:"statusdomain"`
	Created   string       `db:"create_by"`
	CreatedAt sql.NullTime `db:"create_at"`
	Update    string       `db:"update_by"`
	UpdateAt  sql.NullTime `db:"update_at"`
}
type DomainRepository interface {
	FindAll(ctx context.Context) ([]Domain, error)
	FindByID(ctx context.Context, id string) (Domain, error)
	Save(ctx context.Context, cur *Domain) error
	Update(ctx context.Context, cur *Domain) error
}
type DomainService interface {
	All(ctx context.Context) ([]dto.DomainData, error)
	Save(ctx context.Context, req dto.DomainSave, client string) error
}
