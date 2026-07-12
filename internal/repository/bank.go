package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type bankRepository struct {
	db DBExecutor
}

func NewBankRepository(db DBExecutor) domain.BankRepository {
	return &bankRepository{
		db: db,
	}
}
func (b bankRepository) FindAll(ctx context.Context) ([]domain.Bank, error) {
	query := `SELECT * FROM ` + config.DB_tbl_bank + ` ORDER BY idbank ASC`

	rows, err := b.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Bank
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Bank])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b bankRepository) FindSelect(ctx context.Context) ([]domain.Bank, error) {
	query := `SELECT idbank, nmbank FROM ` + config.DB_tbl_bank + ` ORDER BY idbank ASC`

	rows, err := b.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Bank])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b bankRepository) FindByID(ctx context.Context, id string) (domain.Bank, error) {
	var bank domain.Bank
	query := `SELECT idbank FROM ` + config.DB_tbl_bank + ` WHERE idbank = $1 LIMIT 1`

	err := b.db.QueryRow(ctx, query, id).Scan(&bank.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Bank{}, nil
		}
		return bank, err
	}
	return bank, nil
}

func (b bankRepository) Save(ctx context.Context, bank *domain.Bank) error {
	query := `INSERT INTO ` + config.DB_tbl_bank + ` 
                (idbank, typebank, nmbank, bankstatus, createbank, createdatebank) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := b.db.Exec(ctx, query,
		bank.ID,
		bank.Type,
		bank.Name,
		bank.Status,
		bank.Created, // Sesuaikan nama field di struct domain Anda
		bank.CreatedAt,
	)
	return err
}

func (b bankRepository) Update(ctx context.Context, bank *domain.Bank) error {
	query := `UPDATE ` + config.DB_tbl_bank + ` SET 
                typebank = $1, 
                nmbank = $2, 
                bankstatus = $3, 
                updatebank = $4, 
                updatedatebank = $5 
              WHERE idbank = $6`

	res, err := b.db.Exec(ctx, query,
		bank.Type,
		bank.Name,
		bank.Status,
		bank.Update,
		bank.UpdateAt,
		bank.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
