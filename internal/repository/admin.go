package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/config"
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
func (a adminRepository) FindAll(ctx context.Context, idcomp string) ([]domain.Admin, error) {
	query := `SELECT * FROM ` + config.DB_tbl_companyadmin + ` 
			WHERE idcompany = $1 
			ORDER BY GREATEST(createdatecompadmin, lastlogincompadmin) DESC`

	rows, err := a.db.Query(ctx, query, idcomp)
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
                passcompadmin, idcompany, idclientrule, namecompadmin  
              FROM ` + config.DB_tbl_companyadmin + ` 
              WHERE usernamecompadmin=$1 AND compadminstatus='Y' LIMIT 1`

	err := a.db.QueryRow(ctx, query, username).Scan(
		&c.Pass,
		&c.IDCompany,
		&c.IDClientrule,
		&c.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Admin{}, nil
		}
		return c, err
	}
	return c, nil
}
func (a adminRepository) FindByUsernameComp(ctx context.Context, username, idcomp string) (domain.Admin, error) {
	var c domain.Admin
	query := `SELECT 
            	idcompadmin, passcompadmin, idclientrule, namecompadmin  
              FROM ` + config.DB_tbl_companyadmin + ` 
              WHERE usernamecompadmin=$1 
			  AND idcompany=$2 
			  AND compadminstatus='Y' LIMIT 1`

	err := a.db.QueryRow(ctx, query, username, idcomp).Scan(
		&c.ID,
		&c.Pass,
		&c.IDClientrule,
		&c.Name,
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
	query := `INSERT INTO ` + config.DB_tbl_companyadmin + ` 
                (idcompadmin, idcompany, idclientrule, usernamecompadmin, passcompadmin, 
                 namecompadmin, compadminstatus, createcompadmin, createdatecompadmin) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := a.db.Exec(ctx, query,
		admin.ID,
		admin.IDCompany,
		admin.IDClientrule,
		admin.Username,
		admin.Pass,
		admin.Name,
		admin.Status,
		admin.Created,
		admin.CreatedAt,
	)
	return err
}

func (a adminRepository) Update(ctx context.Context, admin *domain.Admin, flagpass bool) error {
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
                  WHERE idcompadmin = $7 AND idcompany=$8`
		args = []any{
			admin.IDClientrule,
			admin.Pass,
			admin.Name,
			admin.Status,
			admin.Update,
			admin.UpdateAt,
			admin.ID,
			admin.IDCompany,
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
			admin.IDClientrule,
			admin.Name,
			admin.Status,
			admin.Update,
			admin.UpdateAt,
			admin.ID,
		}
	}

	res, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (a adminRepository) UpdateLogin(ctx context.Context, admin *domain.Admin) error {
	query := `UPDATE ` + config.DB_tbl_companyadmin + ` SET 
                lastlogincompadmin = $1, 
                ipaddresscompadmin = $2  
              WHERE usernamecompadmin = $3 AND idcompany = $4 `
	res, err := a.db.Exec(ctx, query,
		admin.Lastlogin,
		admin.Ipaddress,
		admin.Username,
		admin.IDCompany,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
