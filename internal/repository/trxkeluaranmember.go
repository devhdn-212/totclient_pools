package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/config"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/jackc/pgx/v5"
)

// trxkeluaranmemberColumns lists exactly the columns
// domain.Trxkeluaranmember has fields for — see trxkeluarandetailColumns for
// why SELECT * isn't used here.
const trxkeluaranmemberColumns = `idkeluaranmember, idtrxkeluaran, idcompany,
	username, totalbet, totalbayar, totaldiscount, totalkei, totalwin,
	totalpair, betround, playerinvoice, status,
	createkeluaranmember, createdatekeluaranmember,
	COALESCE(updatekeluaranmember, '') AS updatekeluaranmember, updatedatekeluaranmember`

type trxkeluaranmemberRepository struct {
	db DBExecutor
}

func NewTrxkeluaranmemberRepository(db DBExecutor) domain.TrxkeluaranmemberRepository {
	return &trxkeluaranmemberRepository{
		db: db,
	}
}
func (a trxkeluaranmemberRepository) FindAll(ctx context.Context, idcomp string, idtrx int) ([]domain.Trxkeluaranmember, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluaranmemberColumns + ` FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			WHERE idtrxkeluaran = $1
			ORDER BY createdatekeluaranmember DESC`

	rows, err := a.db.Query(ctx, query, idtrx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Trxkeluaranmember])
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindPeriodsByUsername aggregates every member row per idtrxkeluaran
// (totalpair/totalbayar/totalwin summed across playerinvoices) instead of
// returning one row per invoice — a player's "Transaksi" period list only
// ever needs one card per period, so grouping happens in SQL rather than
// pulling every raw row and grouping client-side. Joins to keluarantogel
// (idcomppasaran scoping + keluarantogel/datekeluaran/keluaranperiode) and
// tbl_mst_company_pasaran (the pasaran's name/code for the exact period
// being shown, instead of the client assuming it's always whatever pasaran
// happens to be currently active).
//
// keluarantogel (empty = draw not decided yet) is what the caller uses to
// resolve a single pending/complete verdict for the period.
//
// No LIMIT: this result set is the period list itself, so truncating it
// would drop whole periods from the player's "Transaksi" history, not just
// old rows.
func (a trxkeluaranmemberRepository) FindPeriodsByUsername(ctx context.Context, idcomp, idcomppasaran, username string) ([]domain.TrxkeluaranmemberPeriod, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT M.idtrxkeluaran,
			K.datekeluaran, K.keluaranperiode, K.keluarantogel,
			P.aliascomppasaran, P.codecomppasaran,
			SUM(M.totalpair) AS totalpair,
			SUM(M.totalbayar) AS totalbayar,
			SUM(M.totalwin) AS totalwin
			FROM ` + t.Schema + `.` + t.KeluarantogelMember + ` AS M
			INNER JOIN ` + t.Schema + `.` + t.Keluarantogel + ` AS K ON K.idtrxkeluaran = M.idtrxkeluaran
			INNER JOIN ` + config.DB_mst_company_pasaran + ` AS P ON P.idcomppasaran = K.idcomppasaran
			WHERE K.idcompany = $1 AND K.idcomppasaran = $2 AND M.username = $3
			GROUP BY M.idtrxkeluaran, K.datekeluaran, K.keluaranperiode, K.keluarantogel, P.aliascomppasaran, P.codecomppasaran
			ORDER BY K.datekeluaran DESC`

	rows, err := a.db.Query(ctx, query, idcomp, idcomppasaran, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.TrxkeluaranmemberPeriod])
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (a trxkeluaranmemberRepository) FindByID(ctx context.Context, idcomp, username string, idtrx int) (domain.Trxkeluaranmember, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluaranmemberColumns + `
			FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			WHERE username = $1 AND idtrxkeluaran = $2
			LIMIT 1`

	rows, err := a.db.Query(ctx, query, idtrx, username, idtrx)
	if err != nil {
		return domain.Trxkeluaranmember{}, err
	}
	defer rows.Close()

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[domain.Trxkeluaranmember])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluaranmember{}, nil
		}
		return domain.Trxkeluaranmember{}, err
	}
	return data, nil
}

// ExistsByUsername reports whether username already has a member row for
// idtrxkeluaran — checkout checks this before saving its own member row to
// decide whether this invoice is the player's first for the period.
func (a trxkeluaranmemberRepository) ExistsByUsername(ctx context.Context, idcomp string, idtrxkeluaran int, username string) (bool, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT EXISTS(
			SELECT 1 FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			WHERE idtrxkeluaran = $1 AND username = $2)`

	var exists bool
	if err := a.db.QueryRow(ctx, query, idtrxkeluaran, username).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (a trxkeluaranmemberRepository) Save(ctx context.Context, trxkeluaranmember *domain.Trxkeluaranmember, idcomp string) error {
	t := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + t.Schema + `.` + t.KeluarantogelMember + `
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

	t := util.Get_mapping_totodb(idcomp)
	query = `UPDATE ` + t.Schema + `.` + t.KeluarantogelMember + ` SET
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
