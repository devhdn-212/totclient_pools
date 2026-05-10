package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"

	"github.com/jackc/pgx/v5"
)

type domainRepository struct {
	db DBExecutor
}

func NewDomainRepository(db DBExecutor) domain.DomainRepository {
	return &domainRepository{
		db: db,
	}
}
func (d domainRepository) FindAll(ctx context.Context) ([]domain.Domain, error) {
	query := `SELECT * FROM ` + config.DB_tbl_domain + ` ORDER BY createdatedomain ASC`

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan pgx v5
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Domain])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (d domainRepository) FindByID(ctx context.Context, id string) (domain.Domain, error) {
	var dm domain.Domain
	query := `SELECT iddomain FROM ` + config.DB_tbl_domain + ` WHERE iddomain = $1 LIMIT 1`

	err := d.db.QueryRow(ctx, query, id).Scan(&dm.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Domain{}, nil
		}
		return dm, err
	}
	return dm, nil
}

func (d domainRepository) Save(ctx context.Context, dm *domain.Domain) error {
	query := `INSERT INTO ` + config.DB_tbl_domain + ` 
                (iddomain, tipedomain, nmdomain, statusdomain, createdomain, createdatedomain) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := d.db.Exec(ctx, query,
		dm.ID,
		dm.Type,
		dm.Name,
		dm.Status,
		dm.Created,
		dm.CreatedAt,
	)
	return err
}

func (d domainRepository) Update(ctx context.Context, dm *domain.Domain) error {
	query := `UPDATE ` + config.DB_tbl_domain + ` SET 
                tipedomain = $1, 
                nmdomain = $2, 
                statusdomain = $3, 
                updatedomain = $4, 
                updatedatedomain = $5 
              WHERE iddomain = $6`

	res, err := d.db.Exec(ctx, query,
		dm.Type,
		dm.Name,
		dm.Status,
		dm.Update,
		dm.UpdateAt,
		dm.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
