package repository

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type bankRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewBankRepository(exec *GoquExecutor) domain.BankRepository {
	return &bankRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (b bankRepository) FindAll(ctx context.Context) ([]domain.Bank, error) {
	var res []domain.Bank
	err := b.db.
		From(config.DB_tbl_bank).
		Order(goqu.C("idbank").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (b bankRepository) FindSelect(ctx context.Context) ([]domain.Bank, error) {
	var res []domain.Bank
	err := b.db.
		From(config.DB_tbl_bank).
		Select("idbank", "nmbank").
		Order(goqu.C("idbank").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (b bankRepository) FindByID(ctx context.Context, id string) (domain.Bank, error) {
	var bank domain.Bank

	ds := b.db.From(config.DB_tbl_bank).
		Select(goqu.C("idbank")).
		Where(
			goqu.C("idbank").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return bank, err
	}

	row := b.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&bank.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Bank{}, nil
	}
	return bank, err
}

func (b bankRepository) Save(ctx context.Context, bank *domain.Bank) error {
	sqlStr, args, err := b.db.Insert(config.DB_tbl_bank).Rows(bank).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (b bankRepository) Update(ctx context.Context, bank *domain.Bank) error {
	sqlStr, args, err := b.db.
		Update(config.DB_tbl_bank).
		Set(goqu.Record{
			"typebank":       bank.Type,
			"nmbank":         bank.Name,
			"bankstatus":     bank.Status,
			"updatebank":     bank.Update,
			"updatedatebank": bank.UpdateAt,
		}).
		Where(goqu.C("idbank").Eq(bank.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := b.exec.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
