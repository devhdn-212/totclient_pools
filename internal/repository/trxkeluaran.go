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
func (a trxkeluaranRepository) FindByID(ctx context.Context, idcomp, idcomppasaran string) (domain.Trxkeluaran, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT
	        A.idtrxkeluaran, A.idcomppasaran, A.keluaranperiode, A.datekeluaran, A.keluarantogel 
			FROM ` + t.Schema + `.` + t.Keluarantogel + ` as A
			WHERE A.idcompany = $1
			AND A.idcomppasaran = $2
			AND A.keluarantogel = ''
			AND A.revisi = 0
			ORDER BY A.create_at DESC LIMIT 1`

	rows, err := a.db.Query(ctx, query, idcomp, idcomppasaran)
	if err != nil {
		return domain.Trxkeluaran{}, err
	}
	defer rows.Close()

	trxkeluaran, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[domain.Trxkeluaran])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluaran{}, nil
		}
		return domain.Trxkeluaran{}, err
	}
	return trxkeluaran, nil
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
