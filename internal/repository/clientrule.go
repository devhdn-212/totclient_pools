package repository

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type clientruleRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewClientruleRepository(exec *GoquExecutor) domain.ClientruleRepository {
	return &clientruleRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (c clientruleRepository) FindAll(ctx context.Context) ([]domain.Clientrule, error) {
	var res []domain.Clientrule
	err := c.db.
		From(config.DB_tbl_clientrule).
		Order(goqu.C("idclientrule").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (c clientruleRepository) FindSelect(ctx context.Context) ([]domain.Clientrule, error) {
	var res []domain.Clientrule
	err := c.db.
		From(config.DB_tbl_clientrule).
		Select("idclientrule", "nmclientrule").
		Order(goqu.C("idclientrule").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (c clientruleRepository) FindByID(ctx context.Context, id string) (domain.Clientrule, error) {
	var cr domain.Clientrule

	ds := c.db.From(config.DB_tbl_clientrule).
		Select("idclientrule").
		Where(
			goqu.C("idclientrule").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return cr, err
	}

	row := c.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&cr.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Clientrule{}, nil
	}
	return cr, err
}

func (c clientruleRepository) Save(ctx context.Context, clientrule *domain.Clientrule) error {
	sqlStr, args, err := c.db.Insert(config.DB_tbl_clientrule).Rows(clientrule).ToSQL()
	if err != nil {
		return err
	}
	_, err = c.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (c clientruleRepository) Update(ctx context.Context, clientrule *domain.Clientrule) error {
	sqlStr, args, err := c.db.
		Update(config.DB_tbl_clientrule).
		Set(goqu.Record{
			"nmclientrule":         clientrule.Name,
			"ruleclient":           clientrule.Rule,
			"updateclientrule":     clientrule.Update,
			"updatedateclientrule": clientrule.UpdateAt,
		}).
		Where(goqu.C("idclientrule").Eq(clientrule.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := c.exec.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
