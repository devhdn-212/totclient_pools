package repository

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"
	"github.com/doug-martin/goqu/v9"
)

type uomRepository struct {
	db   GoquDB
	exec DBExecutor
}

func NewUomRepository(exec *GoquExecutor) domain.UomRepository {
	return &uomRepository{
		db:   exec.DB,
		exec: exec.Exec,
	}
}

// FindAll implements [domain.UomRepository].
func (u *uomRepository) FindAll(ctx context.Context) ([]domain.Uom, error) {
	var res []domain.Uom
	err := u.db.
		From(config.DB_tbl_uom).
		Order(goqu.C("idcurr").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

// FindByID implements [domain.UomRepository].
func (u *uomRepository) FindByID(ctx context.Context, id string) (domain.Uom, error) {
	var record domain.Uom

	ds := u.db.From(config.DB_tbl_uom).
		Select(goqu.C("iduom")).
		Where(
			goqu.C("iduom").Eq(id),
		)

	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return record, err
	}

	row := u.exec.QueryRowContext(ctx, sqlStr, args...)
	err = row.Scan(
		&record.ID,
	)
	if err == sql.ErrNoRows {
		return domain.Uom{}, nil
	}
	return record, err
}

// FindSelect implements [domain.UomRepository].
func (u *uomRepository) FindSelect(ctx context.Context) ([]domain.Uom, error) {
	var res []domain.Uom
	err := u.db.
		From(config.DB_tbl_uom).
		Select("iduom", "nmuom").
		Order(goqu.C("iduom").Asc()).
		ScanStructsContext(ctx, &res)
	return res, err
}

// Save implements [domain.UomRepository].
func (u *uomRepository) Save(ctx context.Context, uom *domain.Uom) error {
	sqlStr, args, err := u.db.Insert(config.DB_tbl_uom).Rows(uom).ToSQL()
	if err != nil {
		return err
	}
	_, err = u.exec.ExecContext(ctx, sqlStr, args...)
	return err
}

// Update implements [domain.UomRepository].
func (u *uomRepository) Update(ctx context.Context, uom *domain.Uom) error {
	sqlStr, args, err := u.db.
		Update(config.DB_tbl_uom).
		Set(goqu.Record{
			"nmuom":     uom.Name,
			"status":    uom.Status,
			"update_by": uom.Update,
			"update_at": uom.UpdateAt,
		}).
		Where(goqu.C("iduom").Eq(uom.ID)).
		ToSQL()
	if err != nil {
		return err
	}
	res, err := u.exec.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
