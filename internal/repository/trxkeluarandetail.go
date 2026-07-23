package repository

import (
	"context"
	"errors"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/jackc/pgx/v5"
)

// trxkeluarandetailColumns lists exactly the columns domain.Trxkeluarandetail
// has fields for. SELECT * breaks the moment the live table gains a column
// (e.g. detail_uuid) that hasn't been added to the struct yet — pgx's
// RowToStructByName(Lax) only tolerates the struct having extra fields, not
// the row having extra columns — so queries here always name columns
// explicitly instead.
const trxkeluarandetailColumns = `idtrxkeluarandetail, idtrxkeluaran, idcompany,
	datetimedetail, COALESCE(ipaddress, '') AS ipaddress, username, typegame, nomortogel,
	COALESCE(posisitogel, '') AS posisitogel,
	bet, diskon, win, winhasil, cancelbet, kei,
	COALESCE(browsertogel, '') AS browsertogel, COALESCE(devicetogel, '') AS devicetogel,
	statuskeluarandetail, betround, winrev, playerinvoice,
	COALESCE(senddata, '') AS senddata,
	senddatacreatedate, COALESCE(updatedata, '') AS updatedata, updatedatacreatedate,
	create_by, create_at, COALESCE(update_by, '') AS update_by, update_at`

type trxkeluarandetailRepository struct {
	db DBExecutor
}

func NewTrxkeluarandetailRepository(db DBExecutor) domain.TrxkeluarandetailRepository {
	return &trxkeluarandetailRepository{
		db: db,
	}
}
func (a trxkeluarandetailRepository) FindAll(ctx context.Context, idcomp string, idtrx int) ([]domain.Trxkeluarandetail, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluarandetailColumns + ` FROM ` + t.Schema + `.` + t.KeluarantogelDetail + `
			WHERE idtrxkeluaran = $1
			ORDER BY create_at DESC`

	rows, err := a.db.Query(ctx, query, idtrx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Trxkeluarandetail])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a trxkeluarandetailRepository) FindByUsername(ctx context.Context, idcomp string, idtrx int, username string) ([]domain.Trxkeluarandetail, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluarandetailColumns + ` FROM ` + t.Schema + `.` + t.KeluarantogelDetail + `
			WHERE idtrxkeluaran = $1 AND username = $2
			ORDER BY create_at DESC
			LIMIT 100`

	rows, err := a.db.Query(ctx, query, idtrx, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Trxkeluarandetail])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a trxkeluarandetailRepository) FindByID(ctx context.Context, idcomp, idtrxkeluarandetail string, idtrx int) (domain.Trxkeluarandetail, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluarandetailColumns + `
			FROM ` + t.Schema + `.` + t.KeluarantogelDetail + `
			WHERE idtrxkeluarandetail = $1 AND idtrxkeluaran = $2
			LIMIT 1`

	rows, err := a.db.Query(ctx, query, idtrx, idtrxkeluarandetail, idtrx)
	if err != nil {
		return domain.Trxkeluarandetail{}, err
	}
	defer rows.Close()

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[domain.Trxkeluarandetail])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluarandetail{}, nil
		}
		return domain.Trxkeluarandetail{}, err
	}
	return data, nil
}
func (a trxkeluarandetailRepository) SumBet(ctx context.Context, idcomp string, idtrx int, datekeluaran time.Time, username, typegame, nomortogel string) (int64, error) {
	t := util.Get_mapping_totodb(idcomp)

	// datekeluaran is a plain `date` column — pgx hands it back as midnight
	// UTC, not midnight Jakarta. datetimedetail is `timestamptz`, compared by
	// absolute instant, so AddDate()-ing the raw UTC-labeled value and using
	// it directly would shift the boundary by the UTC+7 offset relative to
	// what "5 days before the draw date" actually means in Jakarta. Re-anchor
	// to LocJakarta first — same pattern pasaran.go uses for this exact
	// column — so the comparison lines up with when bets are really placed.
	dayStart := time.Date(datekeluaran.Year(), datekeluaran.Month(), datekeluaran.Day(), 0, 0, 0, 0, util.LocJakarta)

	// Bets against one idtrxkeluaran can land anywhere from a few days before
	// the draw date up through the draw date itself (players bet ahead, and
	// the market can close past midnight) — so the pruning window has to
	// cover that whole range, not just datekeluaran's own day, or it'd
	// silently miss real rows and undercount the reseeded limit. -5/+2 days
	// is a deliberately generous margin around the real betting window.
	from := dayStart.AddDate(0, 0, -5)
	to := dayStart.AddDate(0, 0, 3)

	query := `SELECT COALESCE(SUM(bet), 0) FROM ` + t.Schema + `.` + t.KeluarantogelDetail + `
			WHERE idtrxkeluaran = $1 AND typegame = $2 AND nomortogel = $3
			AND datetimedetail >= $4 AND datetimedetail < $5`
	args := []interface{}{idtrx, typegame, nomortogel, from, to}
	if username != "" {
		query += ` AND username = $6`
		args = append(args, username)
	}

	var total int64
	if err := a.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (a trxkeluarandetailRepository) Save(ctx context.Context, trxkeluarandetail *domain.Trxkeluarandetail, idcomp string) error {
	t := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + t.Schema + `.` + t.KeluarantogelDetail + `
                (idtrxkeluarandetail, idtrxkeluaran, idcompany,
				datetimedetail, ipaddress,
				username, typegame, posisitogel, nomortogel,
				bet,diskon,kei,win, statuskeluarandetail,playerinvoice,betround,
				browsertogel, devicetogel,
				create_by,create_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`

	_, err := a.db.Exec(ctx, query,
		trxkeluarandetail.ID,
		trxkeluarandetail.IDtrxkeluaran,
		trxkeluarandetail.IDcomp,
		trxkeluarandetail.Datekeluarandetail,
		trxkeluarandetail.Ipaddress,
		trxkeluarandetail.Username,
		trxkeluarandetail.Typegame,
		trxkeluarandetail.Posisitogel,
		trxkeluarandetail.Nomortogel,
		trxkeluarandetail.Bet,
		trxkeluarandetail.Diskon,
		trxkeluarandetail.Kei,
		trxkeluarandetail.Win,
		trxkeluarandetail.Statuskeluarandetail,
		trxkeluarandetail.Playerinvoice,
		trxkeluarandetail.Betround,
		trxkeluarandetail.Browsertogel,
		trxkeluarandetail.Devicetogel,
		trxkeluarandetail.Created,
		trxkeluarandetail.CreatedAt,
	)
	return err
}
