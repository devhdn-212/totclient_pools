package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totclient_api/dto"
)

type Adminrule struct {
	ID        string       `db:"idclientrule"`
	Name      string       `db:"nmclientrule"`
	Rule      string       `db:"ruleclient"`
	Created   string       `db:"createclientrule"`
	CreatedAt sql.NullTime `db:"createdateclientrule"`
	Update    string       `db:"updateclientrule"`
	UpdateAt  sql.NullTime `db:"updatedateclientrule"`
}

type AdminruleRepository interface {
	FindAll(ctx context.Context) ([]Adminrule, error)
	FindSelect(ctx context.Context) ([]Adminrule, error)
	FindByID(ctx context.Context, id string) (Adminrule, error)
	GetRule(ctx context.Context, id string) (string, error)
}
type AdminruleService interface {
	All(ctx context.Context) ([]dto.AdminruleData, error)
	Select(ctx context.Context) ([]dto.AdminruleSelect, error)
}
