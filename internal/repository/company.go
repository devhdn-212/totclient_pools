package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"

	"github.com/jackc/pgx/v5"
)

type companyRepository struct {
	db DBExecutor
}

func NewCompanyRepository(db DBExecutor) domain.CompanyRepository {
	return &companyRepository{
		db: db,
	}
}
func (c companyRepository) FindAll(ctx context.Context) ([]domain.Company, error) {
	query := `
		SELECT 
			c.idcompany,
            c.idgroupcomp,
            g.nmgroupcomp, 
            c.idcurrdef,
            c.compname,
            c.endjoin,
            c.amountcomp,
            c.telegramid,
            c.urlapitoto,
            c.urlapislot,
            c.compstatus,
            c.compactivetoto,
            c.compactiveslot,
            c.createcomp,
            c.createdatecomp,
            c.updatecomp,
            c.updatedatecomp
		FROM ` + config.DB_tbl_company + ` c
		LEFT JOIN ` + config.DB_tbl_groupcompany + ` g ON c.idgroupcomp = g.idgroupcomp 
		ORDER BY GREATEST(c.createdatecomp, c.updatedatecomp) DESC
	`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis menggunakan pgx v5
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Company])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c companyRepository) FindByID(ctx context.Context, id string) (domain.Company, error) {
	var company domain.Company
	query := `SELECT idcompany FROM ` + config.DB_tbl_company + ` WHERE idcompany = $1 LIMIT 1`

	err := c.db.QueryRow(ctx, query, id).Scan(&company.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Company{}, nil
		}
		return company, err
	}
	return company, nil
}

func (c companyRepository) Save(ctx context.Context, company *domain.Company) error {
	query := `INSERT INTO ` + config.DB_tbl_company + ` 
                (
					idcompany, idcurrdef, idgroupcomp,
					compname, telegramid, urlapitoto, urlapislot,
					compstatus, compactivetoto, compactiveslot, 
					createcomp, createdatecomp
				) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, $11, $12)`

	_, err := c.db.Exec(ctx, query,
		company.ID,
		company.IDcurrdef,
		company.IDgroupcomp,
		company.Name,
		company.TelegramID,
		company.URLapitoto,
		company.URLapislot,
		company.Status,
		company.Activetoto,
		company.Activeslot,
		company.Created,
		company.CreatedAt,
	)
	return err
}

func (c companyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `UPDATE ` + config.DB_tbl_company + ` SET 
                idcurrdef = $1, 
                idgroupcomp = $2, 
                compname = $3, 
                telegramid = $4, 
                urlapitoto = $5, 
                urlapislot = $6, 
                compstatus = $7, 
                compactivetoto = $8, 
                compactiveslot = $9, 
                updatecomp = $10, 
                updatedatecomp = $11   
              WHERE idcompany = $12`

	res, err := c.db.Exec(ctx, query,
		company.IDcurrdef,
		company.IDgroupcomp,
		company.Name,
		company.TelegramID,
		company.URLapitoto,
		company.URLapislot,
		company.Status,
		company.Activetoto,
		company.Activeslot,
		company.Update,
		company.UpdateAt,
		company.ID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
