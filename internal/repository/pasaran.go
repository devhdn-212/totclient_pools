package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type pasaranRepository struct {
	db DBExecutor
}

func NewPasaranRepository(db DBExecutor) domain.PasaranRepository {
	return &pasaranRepository{
		db: db,
	}
}

func (u *pasaranRepository) FindJadwal(ctx context.Context, idcomppasaran string) ([]domain.Pasaranjadwal, error) {
	query := `SELECT * FROM ` + config.DB_mst_company_jadwaltogel + ` 
			  WHERE idcomppasaran = $1  
			  ORDER BY create_at DESC`

	rows, err := u.db.Query(ctx, query, idcomppasaran)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Pasaranjadwal])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *pasaranRepository) FindByID(ctx context.Context, idcomp, codepasaran string) (domain.Pasaran, error) {
	query := `SELECT *
			FROM ` + config.DB_mst_company_pasaran + `
			WHERE idcompany = $1 AND codecomppasaran=$2 AND statuscompasaran = 'Y'
			LIMIT 1`

	rows, err := u.db.Query(ctx, query, idcomp, codepasaran)
	if err != nil {
		return domain.Pasaran{}, err
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Pasaran])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Pasaran{}, nil
		}
		return domain.Pasaran{}, err
	}
	return record, nil
}
