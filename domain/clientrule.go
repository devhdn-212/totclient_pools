package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totagen_api/dto"
)

type Clientrule struct {
	ID        string       `db:"idclientrule"`
	Name      string       `db:"nmclientrule"`
	Rule      string       `db:"ruleclient"`
	Created   string       `db:"createclientrule"`
	CreatedAt sql.NullTime `db:"createdateclientrule"`
	Update    string       `db:"updateclientrule"`
	UpdateAt  sql.NullTime `db:"updatedateclientrule"`
}
type ClientruleRepository interface {
	FindAll(ctx context.Context) ([]Clientrule, error)
	FindSelect(ctx context.Context) ([]Clientrule, error)
	FindByID(ctx context.Context, id string) (Clientrule, error)
	Save(ctx context.Context, clientrule *Clientrule) error
	Update(ctx context.Context, clientrule *Clientrule) error
}
type ClientruleService interface {
	All(ctx context.Context) ([]dto.ClientruleData, error)
	Select(ctx context.Context) ([]dto.ClientruleSelect, error)
	Save(ctx context.Context, req dto.ClientruleSave, client string) error
}
