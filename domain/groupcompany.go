package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totmaster_api/dto"
)

type Groupcompany struct {
	ID        string       `db:"idgroupcomp"`
	Name      string       `db:"nmgroupcomp"`
	Status    string       `db:"statusgroupcomp"`
	Created   string       `db:"create_by"`
	CreatedAt sql.NullTime `db:"create_at"`
	Update    string       `db:"update_by"`
	UpdateAt  sql.NullTime `db:"update_at"`
}

type GroupcompanyRepository interface {
	FindAll(ctx context.Context) ([]Groupcompany, error)
	FindSelect(ctx context.Context) ([]Groupcompany, error)
	FindByID(ctx context.Context, id string) (Groupcompany, error)
	Save(ctx context.Context, cur *Groupcompany) error
	Update(ctx context.Context, cur *Groupcompany) error
}
type GroupcompanyService interface {
	All(ctx context.Context) ([]dto.GroupcompanyData, error)
	Select(ctx context.Context) ([]dto.GroupcompanySelect, error)
	Save(ctx context.Context, req dto.GroupcompanySave, client string) error
}
