package repository

import (
	"context"
	"database/sql"
	"gofibergocu/domain"
	"gofibergocu/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type adminruleRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewAdminruleRepository(exec *GoquExecutor) domain.AdminruleRepository {
	return &adminruleRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (a adminruleRepository) FindAll(ctx context.Context) ([]domain.Adminrule, error) {
	var res []domain.Adminrule
	err := a.db.
		From(config.DB_tbl_adminrule).
		Order(goqu.C("idadminrole").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}
func (a adminruleRepository) FindSelect(ctx context.Context) ([]domain.Adminrule, error) {
	var res []domain.Adminrule
	err := a.db.
		From(config.DB_tbl_adminrule).
		Select("idadminrole", "nmadminrole").
		Order(goqu.C("idadminrole").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}
func (a adminruleRepository) FindByID(ctx context.Context, id string) (domain.Adminrule, error) {
	var c domain.Adminrule

	ds := a.db.From(config.DB_tbl_adminrule).
		Select("idadminrole").
		Where(
			goqu.C("idadminrole").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return c, err
	}

	row := a.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&c.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Adminrule{}, nil
	}
	return c, err
}
func (a adminruleRepository) GetRule(ctx context.Context, id string) (string, error) {
	var c domain.Adminrule

	ds := a.db.From(config.DB_tbl_adminrule).
		Select("ruleadmin").
		Where(
			goqu.C("idadminrole").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return c.Rule, err
	}

	row := a.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&c.Rule,
	)
	if err == sql.ErrNoRows {
		return c.Rule, nil
	}
	return c.Rule, err
}
func (a adminruleRepository) Save(ctx context.Context, adminrule *domain.Adminrule) error {
	sqlStr, args, err := a.db.Insert(config.DB_tbl_adminrule).Rows(adminrule).ToSQL()
	if err != nil {
		return err
	}
	_, err = a.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (a adminruleRepository) Update(ctx context.Context, adminrule *domain.Adminrule) error {
	sqlStr, args, err := a.db.
		Update(config.DB_tbl_adminrule).
		Set(goqu.Record{
			"nmadminrole":         adminrule.Name,
			"ruleadmin":           adminrule.Rule,
			"updateadminrole":     adminrule.Update,
			"updatedateadminrole": adminrule.UpdateAt,
		}).
		Where(goqu.C("idadminrole").Eq(adminrule.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := a.exec.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
