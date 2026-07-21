package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/config"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/jackc/pgx/v5"
)

type trxkeluaranRepository struct {
	db DBExecutor
}

func NewTrxkeluaranRepository(db DBExecutor) domain.TrxkeluaranRepository {
	return &trxkeluaranRepository{
		db: db,
	}
}
func (a trxkeluaranRepository) FindAllRunning(ctx context.Context, idcomp string) ([]domain.Trxkeluaranview, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT
	        A.idtrxkeluaran, A.idcomppasaran, B.aliascomppasaran as nmpasaran,
			A.keluaranperiode, A.datekeluaran, A.total_member, A.total_bet, A. total_buangan
			FROM ` + t.Schema + `.` + t.Keluarantogel + ` as A
			INNER JOIN ` + config.DB_mst_company_pasaran + ` as B ON B.idcomppasaran = A.idcomppasaran
			WHERE A.idcompany = $1 AND A.keluarantogel = ''
			ORDER BY A.create_at DESC`

	rows, err := a.db.Query(ctx, query, idcomp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Otomatis mapping ke struct domain.Admin
	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Trxkeluaranview])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a trxkeluaranRepository) FindByID(ctx context.Context, idcomp, idcomppasaran string, idtrx int) (domain.Trxkeluaran, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT *
			FROM ` + t.Schema + `.` + t.Keluarantogel + `
			WHERE idtrxkeluaran = $1 AND idcomppasaran = $2
			LIMIT 1`

	rows, err := a.db.Query(ctx, query, idtrx, idcomppasaran)
	if err != nil {
		return domain.Trxkeluaran{}, err
	}
	defer rows.Close()

	compconftoto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Trxkeluaran])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluaran{}, nil
		}
		return domain.Trxkeluaran{}, err
	}
	return compconftoto, nil
}
func (a trxkeluaranRepository) FindByIDByNomorKeluaran(ctx context.Context, idcomp, idcomppasaran string) (domain.Trxkeluaran, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT *
			FROM ` + t.Schema + `.` + t.Keluarantogel + `
			WHERE idcomppasaran = $1 AND keluarantogel = ''
			ORDER BY create_at DESC LIMIT 1`

	rows, err := a.db.Query(ctx, query, idcomppasaran)
	if err != nil {
		return domain.Trxkeluaran{}, err
	}
	defer rows.Close()

	compconftoto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Trxkeluaran])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluaran{}, nil
		}
		return domain.Trxkeluaran{}, err
	}
	return compconftoto, nil
}
func (a trxkeluaranRepository) Save(ctx context.Context, trxkeluaran *domain.Trxkeluaran, idcomp string) error {
	t := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + t.Schema + `.` + t.Keluarantogel + `
                (idtrxkeluaran, idcomppasaran, idcompany, yearmonth, 
				keluaranperiode, datekeluaran, create_by,create_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := a.db.Exec(ctx, query,
		trxkeluaran.ID,
		trxkeluaran.IDcomppasaran,
		trxkeluaran.IDcomp,
		trxkeluaran.Yearmonth,
		trxkeluaran.Keluaranperiode,
		trxkeluaran.Datekeluaran,
		trxkeluaran.Created,
		trxkeluaran.CreatedAt,
	)
	return err
}

func (a trxkeluaranRepository) Update(ctx context.Context, trxkeluaran *domain.Trxkeluaran, idcomp string) error {
	var query string
	var args []any

	t := util.Get_mapping_totodb(idcomp)
	query = `UPDATE ` + t.Schema + `.` + t.Keluarantogel + ` SET
                    keluarantogel = $1, 
                    create_by = $2, 
                    create_at = $3 
                  WHERE idtrxkeluaran = $4 AND idcomppasaran=$5 `
	args = []any{
		trxkeluaran.Keluarantogel,
		trxkeluaran.Update,
		trxkeluaran.UpdateAt,
		trxkeluaran.ID,
		trxkeluaran.IDcomppasaran,
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
