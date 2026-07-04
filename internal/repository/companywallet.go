package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"

	"github.com/jackc/pgx/v5"
)

type companywalletRepository struct {
	db DBExecutor
}

func NewCompanywalletRepository(db DBExecutor) domain.CompanywalletRepository {
	return &companywalletRepository{
		db: db,
	}
}

func (c companywalletRepository) FindAll(ctx context.Context, id string) ([]domain.Companywallet, error) {
	query := `SELECT * FROM ` + config.DB_tbl_companywallet + ` 
              WHERE idcompany = $1 
              ORDER BY amountcompwallet DESC`

	rows, err := c.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan pgx v5
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Companywallet])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c companywalletRepository) FindByID(ctx context.Context, id, idcomp, idcurr string) (domain.Companywallet, error) {
	var compwallet domain.Companywallet
	var query string
	var args []any

	// Logika dinamis diganti ke SQL native
	if id != "" {
		query = `SELECT idcompwallet FROM ` + config.DB_tbl_companywallet + ` 
                 WHERE idcompwallet = $1 AND idcompany = $2 LIMIT 1`
		args = []any{id, idcomp}
	} else {
		query = `SELECT idcompwallet FROM ` + config.DB_tbl_companywallet + ` 
                 WHERE idcompany = $1 AND idcurr = $2 LIMIT 1`
		args = []any{idcomp, idcurr}
	}

	err := c.db.QueryRow(ctx, query, args...).Scan(&compwallet.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Companywallet{}, nil
		}
		return compwallet, err
	}
	return compwallet, nil
}

func (c companywalletRepository) Save(ctx context.Context, companywallet *domain.Companywallet) error {
	query := `INSERT INTO ` + config.DB_tbl_companywallet + ` 
                (idcompwallet, idcompany, idcurr, amountcompwallet, compwalletstatus, createcompwallet, createdatecompwallet) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := c.db.Exec(ctx, query,
		companywallet.ID,
		companywallet.IDcompany,
		companywallet.IDcurr,
		companywallet.Amount, // Pastikan nama field di struct sesuai
		companywallet.Status,
		companywallet.Created,
		companywallet.CreatedAt,
	)
	return err
}

func (c companywalletRepository) Update(ctx context.Context, companywallet *domain.Companywallet) error {
	query := `UPDATE ` + config.DB_tbl_companywallet + ` SET 
                compwalletstatus = $1, 
                updatecompwallet = $2, 
                updatedatecompwallet = $3 
              WHERE idcompwallet = $4 AND idcompany = $5`

	res, err := c.db.Exec(ctx, query,
		companywallet.Status,
		companywallet.Update,
		companywallet.UpdateAt,
		companywallet.ID,
		companywallet.IDcompany,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
