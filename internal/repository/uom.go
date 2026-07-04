package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type uomRepository struct {
	db DBExecutor
}

func NewUomRepository(db DBExecutor) domain.UomRepository {
	return &uomRepository{
		db: db,
	}
}

// FindAll implements [domain.UomRepository].
func (u *uomRepository) FindAll(ctx context.Context) ([]domain.Uom, error) {
	query := `SELECT * FROM ` + config.DB_tbl_uom + ` ORDER BY iduom ASC`

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan pgx v5
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Uom])
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID implements [domain.UomRepository].
func (u *uomRepository) FindByID(ctx context.Context, id string) (domain.Uom, error) {
	var record domain.Uom
	query := `SELECT iduom FROM ` + config.DB_tbl_uom + ` WHERE iduom = $1 LIMIT 1`

	err := u.db.QueryRow(ctx, query, id).Scan(&record.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Uom{}, nil
		}
		return record, err
	}
	return record, nil
}

// FindSelect implements [domain.UomRepository].
func (u *uomRepository) FindSelect(ctx context.Context) ([]domain.Uom, error) {
	query := `SELECT iduom, nmuom FROM ` + config.DB_tbl_uom + ` ORDER BY iduom ASC`

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Uom])
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Save implements [domain.UomRepository].
func (u *uomRepository) Save(ctx context.Context, uom *domain.Uom) error {
	query := `INSERT INTO ` + config.DB_tbl_uom + ` 
                (iduom, nmuom, status, create_by, create_at) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := u.db.Exec(ctx, query,
		uom.ID,
		uom.Name,
		uom.Status,
		uom.Created, // Pastikan field ini sesuai di struct domain.Uom
		uom.CreatedAt,
	)
	return err
}

// Update implements [domain.UomRepository].
func (u *uomRepository) Update(ctx context.Context, uom *domain.Uom) error {
	query := `UPDATE ` + config.DB_tbl_uom + ` SET 
                nmuom = $1, 
                status = $2, 
                update_by = $3, 
                update_at = $4 
              WHERE iduom = $5`

	res, err := u.db.Exec(ctx, query,
		uom.Name,
		uom.Status,
		uom.Update,
		uom.UpdateAt,
		uom.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
