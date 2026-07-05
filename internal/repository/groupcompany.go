package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type groupcompanyRepository struct {
	db DBExecutor
}

func NewGroupcompanyRepository(db DBExecutor) domain.GroupcompanyRepository {
	return &groupcompanyRepository{
		db: db,
	}
}
func (g groupcompanyRepository) FindAll(ctx context.Context) ([]domain.Groupcompany, error) {
	query := `SELECT * FROM ` + config.DB_tbl_groupcompany + ` 
				ORDER BY create_at DESC, update_at DESC`

	rows, err := g.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Groupcompany
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Groupcompany])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g groupcompanyRepository) FindSelect(ctx context.Context) ([]domain.Groupcompany, error) {
	query := `SELECT idgroupcomp, nmgroupcomp FROM ` + config.DB_tbl_groupcompany + ` ORDER BY nmgroupcomp ASC`
	rows, err := g.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Groupcompany])
	if err != nil {
		return nil, err
	}
	// Convert ke domain.Groupcompany
	var result []domain.Groupcompany
	for _, v := range res {
		result = append(result, domain.Groupcompany{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return result, nil
}

func (g groupcompanyRepository) FindByID(ctx context.Context, id string) (domain.Groupcompany, error) {
	var groupcompany domain.Groupcompany
	query := `SELECT idgroupcomp FROM ` + config.DB_tbl_groupcompany + ` WHERE idgroupcomp = $1 LIMIT 1`

	err := g.db.QueryRow(ctx, query, id).Scan(&groupcompany.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Groupcompany{}, nil
		}
		return groupcompany, err
	}
	return groupcompany, nil
}

func (g groupcompanyRepository) Save(ctx context.Context, groupcompany *domain.Groupcompany) error {
	query := `INSERT INTO ` + config.DB_tbl_groupcompany + ` 
                (idgroupcomp, nmgroupcomp, statusgroupcomp, create_by, create_at) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := g.db.Exec(ctx, query,
		groupcompany.ID,
		groupcompany.Name,
		groupcompany.Status,
		groupcompany.Created,
		groupcompany.CreatedAt,
	)
	return err
}

func (g groupcompanyRepository) Update(ctx context.Context, groupcompany *domain.Groupcompany) error {
	query := `UPDATE ` + config.DB_tbl_groupcompany + ` SET 
                nmgroupcomp = $1, 
                statusgroupcomp = $2, 
                update_by = $3, 
                update_at = $4 
              WHERE idgroupcomp = $5`

	res, err := g.db.Exec(ctx, query,
		groupcompany.Name,
		groupcompany.Status,
		groupcompany.Update,
		groupcompany.UpdateAt,
		groupcompany.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
