package repository

import (
	"context"
	"errors"
	"time"

	"github.com/devhdn-212/gofibermaster_api/domain"

	"github.com/jackc/pgx/v5"
)

type customerRepository struct {
	db DBExecutor
}

func NewCustomerRepository(db DBExecutor) domain.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r customerRepository) FindAll(ctx context.Context) ([]domain.Customer, error) {
	query := `SELECT id, code, name, created_at, updated_at, deleted_at 
	          FROM tbl_customer WHERE deleted_at IS NULL`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// pgx/v5 memiliki RowToStructByName yang menggantikan fitur ScanStructs di Goqu
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Customer])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r customerRepository) FindByID(ctx context.Context, id string) (domain.Customer, error) {
	var c domain.Customer
	query := `SELECT id, code, name, created_at, updated_at, deleted_at 
	          FROM tbl_customer WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Code,
		&c.Name,
		&c.CreatedAt,
		&c.UpdateAt, // Pastikan nama field di struct domain sesuai
		&c.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Customer{}, nil
		}
		return c, err
	}
	return c, nil
}

func (r customerRepository) FindByCode(ctx context.Context, code string) (domain.Customer, error) {
	query := `SELECT id, code, name, created_at, updated_at, deleted_at 
	          FROM tbl_customer WHERE code = $1 AND deleted_at IS NULL LIMIT 1`

	rows, err := r.db.Query(ctx, query, code)
	if err != nil {
		return domain.Customer{}, err
	}

	// Menggunakan CollectOneRow untuk mengambil satu data saja secara praktis
	c, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Customer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Customer{}, nil
		}
		return domain.Customer{}, err
	}

	return c, nil
}

func (r customerRepository) Save(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO tbl_customer (id, code, name, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, c.ID, c.Code, c.Name, time.Now(), time.Now())
	return err
}

func (r customerRepository) Update(ctx context.Context, c *domain.Customer) error {
	query := `UPDATE tbl_customer SET name = $1, updated_at = $2 
	          WHERE id = $3 AND deleted_at IS NULL`

	res, err := r.db.Exec(ctx, query, c.Name, time.Now(), c.ID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r customerRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tbl_customer SET deleted_at = $1 WHERE id = $2`

	res, err := r.db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
