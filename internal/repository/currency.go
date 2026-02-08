package repository

import (
	"context"
	"database/sql"
	"gofibergocu/domain"
	"gofibergocu/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type currRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewCurrRepository(exec *GoquExecutor) domain.CurrencyRepository {
	return &currRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (c currRepository) FindAll(ctx context.Context) ([]domain.Currency, error) {
	var res []domain.Currency
	err := c.db.
		From(config.DB_tbl_currency).
		Order(goqu.C("idcurr").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}
func (c currRepository) FindByID(ctx context.Context, id string) (domain.Currency, error) {
	var curr domain.Currency

	ds := c.db.From(config.DB_tbl_currency).
		Select(goqu.C("idcurr")).
		Where(
			goqu.C("idcurr").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return curr, err
	}

	row := c.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&curr.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Currency{}, nil
	}
	return curr, err
}
func (c currRepository) Save(ctx context.Context, cur *domain.Currency) error {
	sqlStr, args, err := c.db.Insert(config.DB_tbl_currency).Rows(cur).ToSQL()
	if err != nil {
		return err
	}
	_, err = c.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (c currRepository) Update(ctx context.Context, cur *domain.Currency) error {
	sqlStr, args, err := c.db.
		Update(config.DB_tbl_currency).
		Set(goqu.Record{
			"typecurr":       cur.Type,
			"status":         cur.Status,
			"updatecurr":     cur.Update,
			"updatedatecurr": cur.UpdateAt,
		}).
		Where(goqu.C("idcurr").Eq(cur.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := c.exec.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
