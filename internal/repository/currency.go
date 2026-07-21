package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/config"

	"github.com/jackc/pgx/v5"
)

type currRepository struct {
	db DBExecutor
}

func NewCurrRepository(db DBExecutor) domain.CurrencyRepository {
	return &currRepository{
		db: db,
	}
}
func (c currRepository) FindAll(ctx context.Context) ([]domain.Currency, error) {
	query := `SELECT * FROM ` + config.DB_tbl_currency + ` ORDER BY idcurr ASC`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan pgx v5
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Currency])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (c currRepository) FindSelect(ctx context.Context) ([]domain.Currency, error) {
	query := `SELECT idcurr FROM ` + config.DB_tbl_currency + ` ORDER BY idcurr ASC`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Currency])
	if err != nil {
		return nil, err
	}

	// Convert ke domain.Groupcompany
	var result []domain.Currency
	for _, v := range res {
		result = append(result, domain.Currency{
			ID: v.ID,
		})
	}

	return res, nil
}
func (c currRepository) FindByID(ctx context.Context, id string) (domain.Currency, error) {
	var curr domain.Currency
	query := `SELECT idcurr FROM ` + config.DB_tbl_currency + ` WHERE idcurr = $1 LIMIT 1`

	err := c.db.QueryRow(ctx, query, id).Scan(&curr.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Currency{}, nil
		}
		return curr, err
	}
	return curr, nil
}
func (c currRepository) Save(ctx context.Context, cur *domain.Currency) error {
	query := `INSERT INTO ` + config.DB_tbl_currency + ` 
                (idcurr, typecurr, statuscurr, createcurr, createdatecurr) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := c.db.Exec(ctx, query,
		cur.ID,
		cur.Type,
		cur.Status,
		cur.Created,
		cur.CreatedAt,
	)
	return err
}

func (c currRepository) Update(ctx context.Context, cur *domain.Currency) error {
	query := `UPDATE ` + config.DB_tbl_currency + ` SET 
                typecurr = $1, 
                statuscurr = $2, 
                updatecurr = $3, 
                updatedatecurr = $4 
              WHERE idcurr = $5`

	res, err := c.db.Exec(ctx, query,
		cur.Type,
		cur.Status,
		cur.Update,
		cur.UpdateAt,
		cur.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
