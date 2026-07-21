package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type companyconftotoRepository struct {
	db DBExecutor
}

func NewCompanyconftotoRepository(db DBExecutor) domain.CompanyconftotoRepository {
	return &companyconftotoRepository{
		db: db,
	}
}

func (c companyconftotoRepository) FindAll(ctx context.Context, idcompany string) ([]domain.Companyconftoto, error) {
	query := `SELECT * FROM ` + config.DB_tbl_companyconftoto + ` 
              WHERE idcompany = $1 
			  ORDER BY GREATEST(create_at, update_at) DESC`

	rows, err := c.db.Query(ctx, query, idcompany)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Companyconftoto])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c companyconftotoRepository) FindByID(ctx context.Context, idcompany string) (domain.Companyconftoto, error) {
	query := `SELECT * FROM ` + config.DB_tbl_companyconftoto + `
              WHERE idcompany = $1 LIMIT 1`

	rows, err := c.db.Query(ctx, query, idcompany)
	if err != nil {
		return domain.Companyconftoto{}, err
	}
	defer rows.Close()

	compconftoto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Companyconftoto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Companyconftoto{}, nil
		}
		return domain.Companyconftoto{}, err
	}
	return compconftoto, nil
}
