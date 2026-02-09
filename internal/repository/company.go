package repository

import (
	"context"
	"database/sql"
	"gofibergocu/domain"
	"gofibergocu/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type companyRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewCompanyRepository(exec *GoquExecutor) domain.CompanyRepository {
	return &companyRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (c companyRepository) FindAll(ctx context.Context) ([]domain.Company, error) {
	var res []domain.Company
	err := c.db.
		From(config.DB_tbl_company).
		Order(goqu.C("idcompany").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (c companyRepository) FindByID(ctx context.Context, id string) (domain.Company, error) {
	var company domain.Company

	ds := c.db.From(config.DB_tbl_company).
		Select(goqu.C("idcompany")).
		Where(
			goqu.C("idcompany").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return company, err
	}

	row := c.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&company.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Company{}, nil
	}
	return company, err
}

func (c companyRepository) Save(ctx context.Context, company *domain.Company) error {
	sqlStr, args, err := c.db.Insert(config.DB_tbl_company).Rows(company).ToSQL()
	if err != nil {
		return err
	}
	_, err = c.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (c companyRepository) Update(ctx context.Context, company *domain.Company) error {
	sqlStr, args, err := c.db.
		Update(config.DB_tbl_company).
		Set(goqu.Record{
			"idcurrdef":      company.IDcurrdef,
			"compname":       company.Name,
			"compstatus":     company.Status,
			"updatecomp":     company.Update,
			"updatedatecomp": company.UpdateAt,
		}).
		Where(goqu.C("idcompany").Eq(company.ID)).
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
