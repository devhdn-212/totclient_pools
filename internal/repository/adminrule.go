package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"
	"github.com/jackc/pgx/v5"
)

type adminruleRepository struct {
	db DBExecutor
}

func NewAdminruleRepository(db DBExecutor) domain.AdminruleRepository {
	return &adminruleRepository{
		db: db,
	}
}
func (a adminruleRepository) FindAll(ctx context.Context) ([]domain.Adminrule, error) {
	query := `SELECT * FROM ` + config.DB_tbl_adminrule + ` ORDER BY idadminrole ASC`

	rows, err := a.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Adminrule])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a adminruleRepository) FindSelect(ctx context.Context) ([]domain.Adminrule, error) {
	query := `SELECT idadminrole, nmadminrole FROM ` + config.DB_tbl_adminrule + ` ORDER BY idadminrole ASC`

	rows, err := a.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping manual per baris
	res, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Adminrule, error) {
		var ar domain.Adminrule
		err := row.Scan(&ar.ID, &ar.Name)
		return ar, err
	})

	return res, nil
}
func (a adminruleRepository) FindByID(ctx context.Context, id string) (domain.Adminrule, error) {
	var c domain.Adminrule
	query := `SELECT idadminrole FROM ` + config.DB_tbl_adminrule + ` WHERE idadminrole = $1 LIMIT 1`

	err := a.db.QueryRow(ctx, query, id).Scan(&c.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Adminrule{}, nil
		}
		return c, err
	}
	return c, nil
}
func (a adminruleRepository) GetRule(ctx context.Context, id string) (string, error) {
	var rule string
	query := `SELECT ruleadmin FROM ` + config.DB_tbl_adminrule + ` WHERE idadminrole = $1 LIMIT 1`

	err := a.db.QueryRow(ctx, query, id).Scan(&rule)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return rule, nil
}
func (a adminruleRepository) Save(ctx context.Context, adminrule *domain.Adminrule) error {
	query := `INSERT INTO ` + config.DB_tbl_adminrule + ` 
                (idadminrole, nmadminrole, ruleadmin, createdadminrole, createddateadminrole) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := a.db.Exec(ctx, query,
		adminrule.ID,
		adminrule.Name,
		adminrule.Rule,
		adminrule.Created,
		adminrule.CreatedAt,
	)
	return err
}

func (a adminruleRepository) Update(ctx context.Context, adminrule *domain.Adminrule) error {
	query := `UPDATE ` + config.DB_tbl_adminrule + ` SET 
                nmadminrole = $1, 
                ruleadmin = $2, 
                updateadminrole = $3, 
                updatedateadminrole = $4 
              WHERE idadminrole = $5`

	res, err := a.db.Exec(ctx, query,
		adminrule.Name,
		adminrule.Rule,
		adminrule.Update,
		adminrule.UpdateAt,
		adminrule.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
