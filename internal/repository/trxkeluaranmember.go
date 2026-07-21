package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/internal/util"
	"github.com/jackc/pgx/v5"
)

type trxkeluaranmemberRepository struct {
	db DBExecutor
}

func NewTrxkeluaranmemberRepository(db DBExecutor) domain.TrxkeluaranmemberRepository {
	return &trxkeluaranmemberRepository{
		db: db,
	}
}
func (a trxkeluaranmemberRepository) FindAll(ctx context.Context, idcomp string, idtrx int) ([]domain.Trxkeluaranmember, error) {
	schema, _, _, _, tbl_trx_keluarantogel_member, _ := util.Get_mapping_totodb(idcomp)
	query := `SELECT * FROM ` + schema + `.` + tbl_trx_keluarantogel_member + ` 
			WHERE idtrxkeluaran = $1 
			ORDER BY createdatekeluaranmember DESC`

	rows, err := a.db.Query(ctx, query, idtrx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Trxkeluaranmember])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a trxkeluaranmemberRepository) FindByID(ctx context.Context, idcomp, username string, idtrx int) (domain.Trxkeluaranmember, error) {
	schema, _, _, _, tbl_trx_keluarantogel_member, _ := util.Get_mapping_totodb(idcomp)
	query := `SELECT * 
			FROM ` + schema + `.` + tbl_trx_keluarantogel_member + ` 
			WHERE username = $1 AND idtrxkeluaran = $2 
			LIMIT 1`

	rows, err := a.db.Query(ctx, query, idtrx, username, idtrx)
	if err != nil {
		return domain.Trxkeluaranmember{}, err
	}
	defer rows.Close()

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Trxkeluaranmember])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluaranmember{}, nil
		}
		return domain.Trxkeluaranmember{}, err
	}
	return data, nil
}
func (a trxkeluaranmemberRepository) Save(ctx context.Context, trxkeluaranmember *domain.Trxkeluaranmember, idcomp string) error {
	schema, _, _, _, tbl_trx_keluarantogel_member, _ := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + schema + `.` + tbl_trx_keluarantogel_member + ` 
                (idkeluaranmember, idtrxkeluaran, idcompany, 
				username, 
				totalbet, totalbayar, totaldiscount, totalkei, totalwin, totalpair, betround,
				playerinvoice, status, 
				createkeluaranmember, createdatekeluaranmember) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13,$14,$15)`

	_, err := a.db.Exec(ctx, query,
		trxkeluaranmember.ID,
		trxkeluaranmember.IDtrxkeluaran,
		trxkeluaranmember.IDcomp,
		trxkeluaranmember.Username,
		trxkeluaranmember.Totalbet,
		trxkeluaranmember.Totalbayar,
		trxkeluaranmember.Totaldiscount,
		trxkeluaranmember.Totalkei,
		trxkeluaranmember.Totalwin,
		trxkeluaranmember.Totalpair,
		trxkeluaranmember.Betround,
		trxkeluaranmember.Playerinvoice,
		trxkeluaranmember.Status,
		trxkeluaranmember.Created,
		trxkeluaranmember.CreatedAt,
	)
	return err
}

func (a trxkeluaranmemberRepository) Update(ctx context.Context, trxkeluaranmember *domain.Trxkeluaranmember, idcomp string) error {
	var query string
	var args []any

	schema, _, _, _, tbl_trx_keluarantogel_member, _ := util.Get_mapping_totodb(idcomp)
	query = `UPDATE ` + schema + `.` + tbl_trx_keluarantogel_member + ` SET 
                    totalwin = $1, 
                    status = $2, 
                    updatekeluaranmember = $3, 
                    updatedatekeluaranmember = $4  
                  WHERE username = $5 AND playerinvoice =$6 AND idtrxkeluaran=$7 `
	args = []any{
		trxkeluaranmember.Totalwin,
		trxkeluaranmember.Status,
		trxkeluaranmember.Update,
		trxkeluaranmember.UpdateAt,
		trxkeluaranmember.Username,
		trxkeluaranmember.Playerinvoice,
		trxkeluaranmember.IDtrxkeluaran,
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
