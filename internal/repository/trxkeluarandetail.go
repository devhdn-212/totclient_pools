package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/internal/util"
	"github.com/jackc/pgx/v5"
)

type trxkeluarandetailRepository struct {
	db DBExecutor
}

func NewTrxkeluarandetailRepository(db DBExecutor) domain.TrxkeluarandetailRepository {
	return &trxkeluarandetailRepository{
		db: db,
	}
}
func (a trxkeluarandetailRepository) FindAll(ctx context.Context, idcomp string, idtrx int) ([]domain.Trxkeluarandetail, error) {
	schema, _, tbl_trx_keluarantogel_detail, _, _, _ := util.Get_mapping_totodb(idcomp)
	query := `SELECT * FROM ` + schema + `.` + tbl_trx_keluarantogel_detail + ` 
			WHERE idtrxkeluaran = $1 
			ORDER BYcreate_at DESC`

	rows, err := a.db.Query(ctx, query, idtrx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Trxkeluarandetail])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a trxkeluarandetailRepository) FindByID(ctx context.Context, idcomp, idtrxkeluarandetail string, idtrx int) (domain.Trxkeluarandetail, error) {
	schema, _, tbl_trx_keluarantogel_detail, _, _, _ := util.Get_mapping_totodb(idcomp)
	query := `SELECT * 
			FROM ` + schema + `.` + tbl_trx_keluarantogel_detail + ` 
			WHERE idtrxkeluarandetail = $1 AND idtrxkeluaran = $2 
			LIMIT 1`

	rows, err := a.db.Query(ctx, query, idtrx, idtrxkeluarandetail, idtrx)
	if err != nil {
		return domain.Trxkeluarandetail{}, err
	}
	defer rows.Close()

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Trxkeluarandetail])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluarandetail{}, nil
		}
		return domain.Trxkeluarandetail{}, err
	}
	return data, nil
}
func (a trxkeluarandetailRepository) Save(ctx context.Context, trxkeluarandetail *domain.Trxkeluarandetail, idcomp string) error {
	schema, _, tbl_trx_keluarantogel_detail, _, _, _ := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + schema + `.` + tbl_trx_keluarantogel_detail + ` 
                (idtrxkeluarandetail, idtrxkeluaran, idcompany, 
				datetimedetail, 
				username, typegame, posisitogel, nomortogel,
				bet,diskon,kei, statuskeluarandetail,playerinvoice,
				create_by,create_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13,$14,$15)`

	_, err := a.db.Exec(ctx, query,
		trxkeluarandetail.ID,
		trxkeluarandetail.IDtrxkeluaran,
		trxkeluarandetail.IDcomp,
		trxkeluarandetail.Datekeluarandetail,
		trxkeluarandetail.Username,
		trxkeluarandetail.Typegame,
		trxkeluarandetail.Posisitogel,
		trxkeluarandetail.Nomortogel,
		trxkeluarandetail.Bet,
		trxkeluarandetail.Diskon,
		trxkeluarandetail.Kei,
		trxkeluarandetail.Statuskeluarandetail,
		trxkeluarandetail.Playerinvoice,
		trxkeluarandetail.Created,
		trxkeluarandetail.CreatedAt,
	)
	return err
}

func (a trxkeluarandetailRepository) Update(ctx context.Context, trxkeluaran *domain.Trxkeluarandetail, idcomp string) error {
	var query string
	var args []any

	schema, _, tbl_trx_keluarantogel_detail, _, _, _ := util.Get_mapping_totodb(idcomp)
	query = `UPDATE ` + schema + `.` + tbl_trx_keluarantogel_detail + ` SET 
                    statuskeluarandetail = $1, 
                    update_by = $2, 
                    update_at = $3 
                  WHERE idtrxkeluarandetail = $4 AND idtrxkeluaran=$5 `
	args = []any{
		trxkeluaran.Statuskeluarandetail,
		trxkeluaran.Update,
		trxkeluaran.UpdateAt,
		trxkeluaran.ID,
		trxkeluaran.IDtrxkeluaran,
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
