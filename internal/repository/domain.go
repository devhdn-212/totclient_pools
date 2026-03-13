package repository

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"

	"github.com/doug-martin/goqu/v9"
)

type domainRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewDomainRepository(exec *GoquExecutor) domain.DomainRepository {
	return &domainRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}
func (d domainRepository) FindAll(ctx context.Context) ([]domain.Domain, error) {
	var res []domain.Domain
	err := d.db.
		From(config.DB_tbl_domain).
		Order(goqu.C("createdatedomain").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}
func (d domainRepository) FindByID(ctx context.Context, id string) (domain.Domain, error) {
	var dm domain.Domain

	ds := d.db.From(config.DB_tbl_domain).
		Select(goqu.C("iddomain")).
		Where(
			goqu.C("iddomain").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return dm, err
	}

	row := d.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&dm.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Domain{}, nil
	}
	return dm, err
}

func (d domainRepository) Save(ctx context.Context, dm *domain.Domain) error {
	sqlStr, args, err := d.db.Insert(config.DB_tbl_domain).Rows(dm).ToSQL()
	if err != nil {
		return err
	}
	_, err = d.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

func (d domainRepository) Update(ctx context.Context, dm *domain.Domain) error {
	sqlStr, args, err := d.db.
		Update(config.DB_tbl_domain).
		Set(goqu.Record{
			"tipedomain":       dm.Type,
			"nmdomain":         dm.Name,
			"statusdomain":     dm.Status,
			"updatedomain":     dm.Update,
			"updatedatedomain": dm.UpdateAt,
		}).
		Where(goqu.C("iddomain").Eq(dm.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := d.exec.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
