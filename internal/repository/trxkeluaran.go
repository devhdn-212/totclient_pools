package repository

import (
	"context"
	"errors"
	"time"

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

// RefreshTotals recomputes total_member/total_bet/total_pairs/total_payout
// for idtrxkeluaran as an absolute value straight from
// tbl_trx_keluarantogel_member — COUNT(DISTINCT username) naturally handles
// "a repeat checkout by the same player doesn't count twice" for free, since
// it doesn't matter how many invoice rows one username has. Self-healing: run
// it twice, or out of order with another period's refresh, and it always
// lands on the same correct numbers — unlike the old "col = col + $delta"
// increment this replaced, which had every concurrent checkout for the same
// period queue up on one row lock (see service.MarkTotalsDirty/FlushDirtyTotals
// for why this is now called from a background ticker instead of inline).
func (a trxkeluaranRepository) RefreshTotals(ctx context.Context, idcomp string, idtrxkeluaran int) error {
	t := util.Get_mapping_totodb(idcomp)
	query := `UPDATE ` + t.Schema + `.` + t.Keluarantogel + ` AS k
			SET total_member = COALESCE(agg.cnt_member, 0),
			    total_bet = COALESCE(agg.sum_bet, 0),
			    total_pairs = COALESCE(agg.sum_pair, 0),
			    total_payout = COALESCE(agg.sum_bayar, 0)
			FROM (
			    SELECT COUNT(DISTINCT username) AS cnt_member,
			           SUM(totalbet) AS sum_bet,
			           SUM(totalpair) AS sum_pair,
			           SUM(totalbayar) AS sum_bayar
			    FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			    WHERE idtrxkeluaran = $1
			) AS agg
			WHERE k.idtrxkeluaran = $1`

	_, err := a.db.Exec(ctx, query, idtrxkeluaran)
	return err
}

// FindResultsByMonth returns every decided period (keluarantogel != ”)
// with datekeluaran in [start, end) for idcomppasaran — the "Result"
// menu's past-draw-numbers list, newest first.
func (a trxkeluaranRepository) FindResultsByMonth(ctx context.Context, idcomp, idcomppasaran string, start, end time.Time) ([]domain.TrxkeluaranResultRow, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT A.idtrxkeluaran, A.keluaranperiode, A.datekeluaran, A.keluarantogel, B.aliascomppasaran
			FROM ` + t.Schema + `.` + t.Keluarantogel + ` AS A
			INNER JOIN ` + config.DB_mst_company_pasaran + ` AS B ON B.idcomppasaran = A.idcomppasaran
			WHERE A.idcompany = $1 AND A.idcomppasaran = $2 AND A.keluarantogel != ''
			AND A.datekeluaran >= $3 AND A.datekeluaran < $4
			ORDER BY A.datekeluaran DESC`

	rows, err := a.db.Query(ctx, query, idcomp, idcomppasaran, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.TrxkeluaranResultRow])
	if err != nil {
		return nil, err
	}
	return res, nil
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
