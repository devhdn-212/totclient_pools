package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/internal/config"
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
	query := `SELECT * FROM ` + config.DB_tbl_clientrule + ` ORDER BY nmclientrule ASC`

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
	query := `SELECT idclientrule, nmclientrule FROM ` + config.DB_tbl_clientrule + ` ORDER BY nmclientrule ASC`

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
	query := `SELECT idclientrule FROM ` + config.DB_tbl_clientrule + ` WHERE idclientrule = $1 LIMIT 1`

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
	query := `SELECT ruleclient FROM ` + config.DB_tbl_clientrule + ` WHERE idclientrule = $1 LIMIT 1`

	err := a.db.QueryRow(ctx, query, id).Scan(&rule)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return rule, nil
}
