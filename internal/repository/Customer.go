package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/domain"

	"github.com/doug-martin/goqu/v9"
)

type customerRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewCustomerRepository(exec *GoquExecutor) domain.CustomerRepository {
	return &customerRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}

func (r customerRepository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	var res []domain.Customer
	err := r.db.
		From("tbl_customer").
		Where(goqu.C("deleted_at").IsNull()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (r customerRepository) FindByID(ctx context.Context, id string) (domain.Customer, error) {
	var c domain.Customer

	ds := r.db.From("tbl_customer").
		Where(
			goqu.C("id").Eq(id),
			goqu.C("deleted_at").IsNull(),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return c, err
	}

	row := r.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&c.ID,
		&c.Code,
		&c.Name,
		&c.CreatedAt,
		&c.UpdateAt,
		&c.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return domain.Customer{}, nil
	}
	return c, err
}

func (r customerRepository) FindByCode(ctx context.Context, code string) (domain.Customer, error) {
	var c domain.Customer
	_, err := r.db.
		From("tbl_customer").
		Where(
			goqu.C("code").Eq(code),
			goqu.C("deleted_at").IsNull(),
		).
		ScanStructContext(ctx, &c)
	return c, err
}

func (r customerRepository) Save(ctx context.Context, c *domain.Customer) error {
	sqlStr, args, err := r.db.Insert("tbl_customer").Rows(c).ToSQL()
	if err != nil {
		return err
	}
	_, err = r.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (r customerRepository) Update(ctx context.Context, c *domain.Customer) error {
	sqlStr, args, err := r.db.
		Update("tbl_customer").
		Set(goqu.Record{
			"name":       c.Name,
			"updated_at": time.Now(),
		}).
		Where(goqu.C("id").Eq(c.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := r.exec.ExecContext(ctx, sqlStr, args...)
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r customerRepository) Delete(ctx context.Context, id string) error {
	sqlStr, args, err := r.db.
		Update("tbl_customer").
		Set(goqu.Record{"deleted_at": time.Now()}).
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := r.exec.ExecContext(ctx, sqlStr, args...)
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
