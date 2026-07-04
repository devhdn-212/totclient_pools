package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"

	"github.com/jackc/pgx/v5"
)

type companyRepository struct {
	db DBExecutor
}

func NewCompanyRepository(db DBExecutor) domain.CompanyRepository {
	return &companyRepository{
		db: db,
	}
}
func (c companyRepository) FindAll(ctx context.Context) ([]domain.Company, error) {
	query := `SELECT * FROM ` + config.DB_tbl_company + ` ORDER BY idcompany ASC`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan pgx v5
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Company])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c companyRepository) FindByID(ctx context.Context, id string) (domain.Company, error) {
	var company domain.Company
	query := `SELECT idcompany FROM ` + config.DB_tbl_company + ` WHERE idcompany = $1 LIMIT 1`

	err := c.db.QueryRow(ctx, query, id).Scan(&company.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Company{}, nil
		}
		return company, err
	}
	return company, nil
}

func (c companyRepository) Save(ctx context.Context, company *domain.Company) error {
	query := `INSERT INTO ` + config.DB_tbl_company + ` 
                (idcompany, idcurrdef, compname, compstatus, createcomp, createdatecomp) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := c.db.Exec(ctx, query,
		company.ID,
		company.IDcurrdef,
		company.Name,
		company.Status,
		company.Created, // Pastikan field ini ada di struct domain.Company
		company.CreatedAt,
	)
	return err
}

func (c companyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `UPDATE ` + config.DB_tbl_company + ` SET 
                idcurrdef = $1, 
                compname = $2, 
                compstatus = $3, 
                updatecomp = $4, 
                updatedatecomp = $5 
              WHERE idcompany = $6`

	res, err := c.db.Exec(ctx, query,
		company.IDcurrdef,
		company.Name,
		company.Status,
		company.Update,
		company.UpdateAt,
		company.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
