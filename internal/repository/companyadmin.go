package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"
	"github.com/jackc/pgx/v5"
)

type companyadminRepository struct {
	db DBExecutor
}

func NewCompanyadminRepository(db DBExecutor) domain.CompanyadminRepository {
	return &companyadminRepository{
		db: db,
	}
}

func (c companyadminRepository) FindAll(ctx context.Context, idcompany string) ([]domain.Companyadmin, error) {
	query := `SELECT * FROM ` + config.DB_tbl_companyadmin + ` 
              WHERE idcompany = $1 
              ORDER BY lastlogincompadmin DESC`

	rows, err := c.db.Query(ctx, query, idcompany)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyadmin
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Companyadmin])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c companyadminRepository) FindByID(ctx context.Context, idcompany, username string) (domain.Companyadmin, error) {
	var compadmin domain.Companyadmin
	query := `SELECT idcompadmin FROM ` + config.DB_tbl_companyadmin + ` 
              WHERE idcompany = $1 AND usernamecompadmin = $2 LIMIT 1`

	err := c.db.QueryRow(ctx, query, idcompany, username).Scan(&compadmin.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Companyadmin{}, nil
		}
		return compadmin, err
	}
	return compadmin, nil
}

func (c companyadminRepository) Save(ctx context.Context, compadmin *domain.Companyadmin) error {
	query := `INSERT INTO ` + config.DB_tbl_companyadmin + ` 
                (idcompadmin, idcompany, idclientrule, usernamecompadmin, passcompadmin, 
                 namecompadmin, compadminstatus, createcompadmin, createdatecompadmin) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := c.db.Exec(ctx, query,
		compadmin.ID,
		compadmin.IDCompany,
		compadmin.IDClientrule,
		compadmin.Username,
		compadmin.Pass,
		compadmin.Name,
		compadmin.Status,
		compadmin.Created,
		compadmin.CreatedAt,
	)
	return err
}

func (c companyadminRepository) Update(ctx context.Context, compadmin *domain.Companyadmin, flagpass bool) error {
	var query string
	var args []any

	if flagpass {
		query = `UPDATE ` + config.DB_tbl_companyadmin + ` SET 
                    idclientrule = $1, 
                    passcompadmin = $2, 
                    namecompadmin = $3, 
                    compadminstatus = $4, 
                    updatecompadmin = $5, 
                    updatedatecompadmin = $6 
                  WHERE idcompadmin = $7`
		args = []any{
			compadmin.IDClientrule,
			compadmin.Pass,
			compadmin.Name,
			compadmin.Status,
			compadmin.Update,
			compadmin.UpdateAt,
			compadmin.ID,
		}
	} else {
		query = `UPDATE ` + config.DB_tbl_companyadmin + ` SET 
                    idclientrule = $1, 
                    namecompadmin = $2, 
                    compadminstatus = $3, 
                    updatecompadmin = $4, 
                    updatedatecompadmin = $5 
                  WHERE idcompadmin = $6`
		args = []any{
			compadmin.IDClientrule,
			compadmin.Name,
			compadmin.Status,
			compadmin.Update,
			compadmin.UpdateAt,
			compadmin.ID,
		}
	}

	res, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
