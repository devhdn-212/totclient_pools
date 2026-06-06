package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibermaster_api/dto"
)

type Uom struct {
	ID        string       `db:"iduom"`
	Name      string       `db:"nmuom"`
	Status    string       `db:"status"`
	Created   string       `db:"create_by"`
	CreatedAt sql.NullTime `db:"create_at"`
	Update    string       `db:"update_by"`
	UpdateAt  sql.NullTime `db:"update_at"`
}

type UomRepository interface {
	FindAll(ctx context.Context) ([]Uom, error)
	FindSelect(ctx context.Context) ([]Uom, error)
	FindByID(ctx context.Context, id string) (Uom, error)
	Save(ctx context.Context, uom *Uom) error
	Update(ctx context.Context, uom *Uom) error
}
type UomService interface {
	All(ctx context.Context) ([]dto.UomData, error)
	Select(ctx context.Context) ([]dto.UomSelect, error)
	Save(ctx context.Context, req dto.UomSave, client string) error
}
