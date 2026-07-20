package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/internal/config"

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
	query := `
		SELECT 
			c.idcompany,
            c.idgroupcomp,
            g.nmgroupcomp, 
            c.idcurrdef,
            c.compname,
            c.endjoin,
            c.amountcomp,
            c.telegramid,
            c.urlapitoto,
            c.urlapislot,
            c.compstatus,
            c.compactivetoto,
            c.compactiveslot,
            c.createcomp,
            c.createdatecomp,
            c.updatecomp,
            c.updatedatecomp
		FROM ` + config.DB_tbl_company + ` c
		LEFT JOIN ` + config.DB_tbl_groupcompany + ` g ON c.idgroupcomp = g.idgroupcomp 
		ORDER BY GREATEST(c.createdatecomp, c.updatedatecomp) DESC
	`

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
