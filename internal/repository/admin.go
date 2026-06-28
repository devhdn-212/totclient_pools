package repository

import (
	"context"
	"errors"
	"time"

	"github.com/devhdn-212/gofibermaster_api/domain"
	"github.com/devhdn-212/gofibermaster_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type adminRepository struct {
	db DBExecutor
}

func NewAdminRepository(db DBExecutor) domain.AdminsRepository {
	return &adminRepository{
		db: db,
	}
}
func (a adminRepository) FindAll(ctx context.Context) ([]domain.Admin, error) {
	// Menggunakan ORDER BY sesuai logika awal
	query := `SELECT * FROM ` + config.DB_tbl_admin + ` ORDER BY lastlogin DESC`

	rows, err := a.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Otomatis mapping ke struct domain.Admin
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Admin])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a adminRepository) FindByUsername(ctx context.Context, username string) (domain.Admin, error) {
	var c domain.Admin
	query := `SELECT 
                username, password, idadmin, name, statuslogin, 
                lastlogin, joindate, ipaddress, timezone, 
                createadmin, createdateadmin, updateadmin, updatedateadmin 
              FROM ` + config.DB_tbl_admin + ` 
              WHERE username = $1 LIMIT 1`

	err := a.db.QueryRow(ctx, query, username).Scan(
		&c.Username,
		&c.Pass,
		&c.Idadmin,
		&c.Name,
		&c.Status,
		&c.Lastlogin,
		&c.Joindate,
		&c.Ipaddress,
		&c.Timezone,
		&c.Created,
		&c.CreatedAt,
		&c.Update,
		&c.UpdateAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Admin{}, nil
		}
		return c, err
	}
	return c, nil
}

func (a adminRepository) Save(ctx context.Context, admin *domain.Admin) error {
	// Gunakan mapping manual atau pastikan urutan kolom sesuai
	query := `INSERT INTO ` + config.DB_tbl_admin + ` 
                (username, password, idadmin, name, statuslogin, ipaddress, createadmin, createdateadmin) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := a.db.Exec(ctx, query,
		admin.Username, admin.Pass, admin.Idadmin, admin.Name,
		admin.Status, admin.Ipaddress, admin.Username, time.Now(),
	)
	return err
}

func (a adminRepository) Update(ctx context.Context, admin *domain.Admin) error {
	query := `UPDATE ` + config.DB_tbl_admin + ` SET 
                password = $1, 
                idadmin = $2, 
                name = $3, 
                statuslogin = $4, 
                ipaddress = $5, 
                updateadmin = $6, 
                updatedateadmin = $7 
              WHERE username = $8`

	res, err := a.db.Exec(ctx, query,
		admin.Pass,
		admin.Idadmin,
		admin.Name,
		admin.Status,
		admin.Ipaddress,
		admin.Username,
		admin.UpdateAt,
		admin.Username,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (a adminRepository) UpdateNotPassword(ctx context.Context, admin *domain.Admin) error {
	query := `UPDATE ` + config.DB_tbl_admin + ` SET 
                idadmin = $1, 
                name = $2, 
                statuslogin = $3, 
                ipaddress = $4, 
                updateadmin = $5, 
                updatedateadmin = $6 
              WHERE username = $7`

	res, err := a.db.Exec(ctx, query,
		admin.Idadmin,
		admin.Name,
		admin.Status,
		admin.Ipaddress,
		admin.Username,
		admin.UpdateAt,
		admin.Username,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (a adminRepository) UpdateLogin(ctx context.Context, admin *domain.Admin) error {
	query := `UPDATE ` + config.DB_tbl_admin + ` SET 
                lastlogin = $1, 
                ipaddress = $2  
              WHERE username = $3`
	res, err := a.db.Exec(ctx, query,
		admin.Lastlogin,
		admin.Ipaddress,
		admin.Username,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
