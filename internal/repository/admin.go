package repository

import (
	"context"
	"database/sql"
	"gofibergocu/domain"
	"gofibergocu/internal/config"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type adminRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewAdminRepository(exec *GoquExecutor) domain.AdminsRepository {
	return &adminRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (a adminRepository) FindAll(ctx context.Context) ([]domain.Admin, error) {
	var res []domain.Admin
	err := a.db.
		From(config.DB_tbl_admin).
		Order(goqu.C("lastlogin").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (a adminRepository) FindByUsername(ctx context.Context, username string) (domain.Admin, error) {
	var c domain.Admin

	ds := a.db.From(config.DB_tbl_admin).
		Where(
			goqu.C("username").Eq(username),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return c, err
	}

	row := a.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&c.Username,
		&c.Pass,
		&c.Idadmin,
		&c.Name,
		&c.Status,
		&c.Lastlogin,
		&c.Joindate,
		&c.Ipaddress,
		&c.Timezone,
		&c.Created,
		&c.CreatedAt,
		&c.Update,
		&c.UpdateAt,
	)
	if err == sql.ErrNoRows {
		return domain.Admin{}, nil
	}
	return c, err
}

func (a adminRepository) Save(ctx context.Context, admin *domain.Admin) error {
	sqlStr, args, err := a.db.Insert(config.DB_tbl_admin).Rows(admin).ToSQL()
	if err != nil {
		return err
	}
	_, err = a.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (a adminRepository) Update(ctx context.Context, admin *domain.Admin) error {
	sqlStr, args, err := a.db.
		Update(config.DB_tbl_admin).
		Set(goqu.Record{
			"password":        admin.Pass,
			"idadmin":         admin.Idadmin,
			"name":            admin.Name,
			"statuslogin":     admin.Status,
			"ipaddress":       admin.Ipaddress,
			"updateadmin":     admin.Username,
			"updatedateadmin": admin.UpdateAt,
		}).
		Where(goqu.C("username").Eq(admin.Username)).
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
func (a adminRepository) UpdateNoPassword(ctx context.Context, admin *domain.Admin) error {
	sqlStr, args, err := a.db.
		Update(config.DB_tbl_admin).
		Set(goqu.Record{
			"idadmin":         admin.Idadmin,
			"name":            admin.Name,
			"statuslogin":     admin.Status,
			"ipaddress":       admin.Ipaddress,
			"update":          admin.Username,
			"updatedateadmin": time.Now(),
		}).
		Where(goqu.C("username").Eq(admin.Username)).
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
