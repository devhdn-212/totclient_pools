package domain

import (
	"context"
	"database/sql"
	"gofibergocu/dto"
)

type Adminrule struct {
	ID        string       `db:"idadminrole"`
	Name      string       `db:"nmadminrole"`
	Rule      string       `db:"ruleadmin"`
	Created   string       `db:"createadminrole"`
	CreatedAt sql.NullTime `db:"createdateadminrole"`
	Update    string       `db:"updateadminrole"`
	UpdateAt  sql.NullTime `db:"updatedateadminrole"`
}

type AdminruleRepository interface {
	FindAll(ctx context.Context) ([]Adminrule, error)
	FindSelect(ctx context.Context) ([]Adminrule, error)
	FindByID(ctx context.Context, id string) (Adminrule, error)
	Save(ctx context.Context, adminrule *Adminrule) error
	Update(ctx context.Context, adminrule *Adminrule) error
}
type AdminruleService interface {
	All(ctx context.Context) ([]dto.AdminruleData, error)
	Select(ctx context.Context) ([]dto.AdminruleSelect, error)
	Save(ctx context.Context, req dto.AdminruleSave, client string) error
}
