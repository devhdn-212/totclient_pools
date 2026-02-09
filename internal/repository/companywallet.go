package repository

import (
	"context"
	"database/sql"
	"gofibergocu/domain"
	"gofibergocu/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type companywalletRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewCompanywalletRepository(exec *GoquExecutor) domain.CompanywalletRepository {
	return &companywalletRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}

func (c companywalletRepository) FindAll(ctx context.Context, id string) ([]domain.Companywallet, error) {
	var res []domain.Companywallet
	err := c.db.From(config.DB_tbl_companywallet).
		Where(
			goqu.C("idcompany").Eq(id),
		).
		Order(goqu.C("amountcompwallet").Desc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (c companywalletRepository) FindByID(ctx context.Context, id, idcomp, idcurr string) (domain.Companywallet, error) {
	var compwallet domain.Companywallet

	ds := c.db.From(config.DB_tbl_companywallet).
		Select(goqu.C("idcompwallet"))
	if id != "" {
		ds = ds.Where(
			goqu.C("idcompwallet").Eq(id),
			goqu.C("idcompany").Eq(idcomp),
		)
	} else {
		ds = ds.Where(
			goqu.C("idcompany").Eq(idcomp),
			goqu.C("idcurr").Eq(idcurr),
		)
	}

	sqlStr, args, err := ds.ToSQL()
	/*log.Printf("SQL: %\ns", sqlStr)
	log.Printf("ARGS: %+v\n", args)*/
	if err != nil {
		return compwallet, err
	}

	row := c.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&compwallet.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Companywallet{}, nil
	}
	return compwallet, err
}

func (c companywalletRepository) Save(ctx context.Context, companywallet *domain.Companywallet) error {
	sqlStr, args, err := c.db.Insert(config.DB_tbl_companywallet).Rows(companywallet).ToSQL()
	if err != nil {
		return err
	}
	_, err = c.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (c companywalletRepository) Update(ctx context.Context, companywallet *domain.Companywallet) error {
	sqlStr, args, err := c.db.
		Update(config.DB_tbl_companywallet).
		Set(goqu.Record{
			"compwalletstatus":     companywallet.Status,
			"updatecompwallet":     companywallet.Update,
			"updatedatecompwallet": companywallet.UpdateAt,
		}).
		Where(
			goqu.C("idcompwallet").Eq(companywallet.ID),
			goqu.C("idcompany").Eq(companywallet.IDcompany),
		).
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
