package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/internal/config"
	"github.com/jackc/pgx/v5"
)

type settingRepository struct {
	db DBExecutor
}

func NewSettingRepository(db DBExecutor) domain.SettingRepository {
	return &settingRepository{
		db: db,
	}
}
func (b settingRepository) FindAll(ctx context.Context) ([]domain.Setting, error) {
	query := `SELECT appversion,startmaintenance,endmaintenance,shio_parent 
				FROM ` + config.DB_tbl_mst_setting + `
				ORDER BY idsetting ASC LIMIT 1`

	rows, err := b.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Setting])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b settingRepository) FindByID(ctx context.Context) (domain.Setting, error) {
	var setting domain.Setting
	query := `SELECT appversion,startmaintenance,endmaintenance,shio_parent
			FROM ` + config.DB_tbl_mst_setting + ` WHERE idsetting = 1 LIMIT 1`

	err := b.db.QueryRow(ctx, query).Scan(
		&setting.Appversion,
		&setting.Startmaintenance,
		&setting.Endmaintenance,
		&setting.Shio_parent,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Setting{}, nil
		}
		return setting, err
	}
	return setting, nil
}
