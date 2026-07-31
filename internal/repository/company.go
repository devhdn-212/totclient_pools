package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/internal/config"
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

func (u *companyRepository) FindByID(ctx context.Context, idcompany string) (domain.Company, error) {
	query := `SELECT idcompany, compname, compstatus, urlapitoto
			FROM ` + config.DB_tbl_company + `
			WHERE idcompany = $1 AND compstatus = 'Y'
			LIMIT 1`

	rows, err := u.db.Query(ctx, query, idcompany)
	if err != nil {
		return domain.Company{}, err
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Company])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Company{}, nil
		}
		return domain.Company{}, err
	}
	return record, nil
}
