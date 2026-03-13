package repository

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type companyadminRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewCompanyadminRepository(exec *GoquExecutor) domain.CompanyadminRepository {
	return &companyadminRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}

func (c companyadminRepository) FindAll(ctx context.Context, idcompany string) ([]domain.Companyadmin, error) {
	var res []domain.Companyadmin
	err := c.db.
		From(config.DB_tbl_companyadmin).
		Where(
			goqu.C("idcompany").Eq(idcompany),
		).
		Order(goqu.C("lastlogincompadmin").Desc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

func (c companyadminRepository) FindByID(ctx context.Context, idcompany, username string) (domain.Companyadmin, error) {
	var compadmin domain.Companyadmin

	ds := c.db.From(config.DB_tbl_companyadmin).
		Select(goqu.C("idcompadmin")).
		Where(
			goqu.C("idcompany").Eq(idcompany),
			goqu.C("usernamecompadmin").Eq(username),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return compadmin, err
	}

	row := c.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&compadmin.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Companyadmin{}, nil
	}
	return compadmin, err
}

func (c companyadminRepository) Save(ctx context.Context, compadmin *domain.Companyadmin) error {
	sqlStr, args, err := c.db.Insert(config.DB_tbl_companyadmin).Rows(compadmin).ToSQL()
	if err != nil {
		return err
	}
	_, err = c.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (c companyadminRepository) Update(ctx context.Context, compadmin *domain.Companyadmin, flagpass bool) error {
	ds := c.db.Update(config.DB_tbl_companyadmin)

	if flagpass {
		ds = ds.Set(goqu.Record{
			"idclientrule":        compadmin.IDClientrule,
			"passcompadmin":       compadmin.Pass,
			"namecompadmin":       compadmin.Name,
			"compadminstatus":     compadmin.Status,
			"updatecompadmin":     compadmin.Update,
			"updatedatecompadmin": compadmin.UpdateAt,
		})
	} else {
		ds = ds.Set(goqu.Record{
			"idclientrule":        compadmin.IDClientrule,
			"namecompadmin":       compadmin.Name,
			"compadminstatus":     compadmin.Status,
			"updatecompadmin":     compadmin.Update,
			"updatedatecompadmin": compadmin.UpdateAt,
		})
	}

	sqlStr, args, err := ds.
		Where(goqu.C("idcompadmin").Eq(compadmin.ID)).
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
