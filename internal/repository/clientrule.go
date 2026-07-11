package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"

	"github.com/jackc/pgx/v5"
)

type clientruleRepository struct {
	db DBExecutor
}

func NewClientruleRepository(db DBExecutor) domain.ClientruleRepository {
	return &clientruleRepository{
		db: db,
	}
}
func (c clientruleRepository) FindAll(ctx context.Context) ([]domain.Clientrule, error) {
	query := `SELECT * FROM ` + config.DB_tbl_clientrule + ` 
				ORDER BY updatedateclientrule DESC`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Clientrule
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Clientrule])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c clientruleRepository) FindSelect(ctx context.Context) ([]domain.Clientrule, error) {
	query := `SELECT idclientrule, nmclientrule FROM ` + config.DB_tbl_clientrule + `
				 ORDER BY idclientrule ASC`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping manual per baris
	res, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Clientrule, error) {
		var cr domain.Clientrule
		err := row.Scan(&cr.ID, &cr.Name)
		return cr, err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c clientruleRepository) FindByID(ctx context.Context, id string) (domain.Clientrule, error) {
	var cr domain.Clientrule
	query := `SELECT idclientrule FROM ` + config.DB_tbl_clientrule + ` 
			WHERE idclientrule = $1  
			LIMIT 1`

	err := c.db.QueryRow(ctx, query, id).Scan(&cr.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Clientrule{}, nil
		}
		return cr, err
	}
	return cr, nil
}

func (c clientruleRepository) Save(ctx context.Context, clientrule *domain.Clientrule) error {
	query := `INSERT INTO ` + config.DB_tbl_clientrule + ` 
                (idclientrule, nmclientrule, ruleclient, createclientrule, createdateclientrule) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := c.db.Exec(ctx, query,
		clientrule.ID,
		clientrule.Name,
		clientrule.Rule,
		clientrule.Created, // Sesuaikan dengan field di struct domain Anda
		clientrule.CreatedAt,
	)
	return err
}

func (c clientruleRepository) Update(ctx context.Context, clientrule *domain.Clientrule) error {
	query := `UPDATE ` + config.DB_tbl_clientrule + ` SET 
                nmclientrule = $1, 
                ruleclient = $2, 
                updateclientrule = $3, 
                updatedateclientrule = $4 
              WHERE idclientrule = $5 `

	res, err := c.db.Exec(ctx, query,
		clientrule.Name,
		clientrule.Rule,
		clientrule.Update,
		clientrule.UpdateAt,
		clientrule.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
