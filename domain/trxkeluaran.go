package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/devhdn-212/totclient_api/dto"
	"github.com/shopspring/decimal"
)

type Trxkeluaran struct {
	ID                int             `db:"idtrxkeluaran"`
	IDcomppasaran     string          `db:"idcomppasaran"`
	IDcomp            string          `db:"idcompany"`
	Yearmonth         string          `db:"yearmonth"`
	Keluaranperiode   int             `db:"keluaranperiode"`
	Datekeluaran      time.Time       `db:"datekeluaran"`
	Keluarantogel     string          `db:"keluarantogel"`
	Prize2            string          `db:"prize2"`
	Prize3            string          `db:"prize3"`
	Total_member      int             `db:"total_member"`
	Total_bet         decimal.Decimal `db:"total_bet"`
	Total_outstanding decimal.Decimal `db:"total_pairs"`
	Total_win         decimal.Decimal `db:"total_win"`
	Total_lose        decimal.Decimal `db:"total_lose"`
	Total_buangan     decimal.Decimal `db:"total_payout"`
	Total_reject      decimal.Decimal `db:"total_reject"`
	Total_winlose     decimal.Decimal `db:"winlose"`
	Total_revisi      int             `db:"revisi"`
	Noterevisi        string          `db:"noterevisi"`
	Created           string          `db:"create_by"`
	CreatedAt         sql.NullTime    `db:"create_at"`
	Update            string          `db:"update_by"`
	UpdateAt          sql.NullTime    `db:"update_at"`
}
type Trxkeluaranview struct {
	ID                int             `db:"idtrxkeluaran"`
	IDcomppasaran     string          `db:"idcomppasaran"`
	IDcomp            string          `db:"idcompany"`
	Nmpasaran         string          `db:"nmpasaran"`
	Keluaranperiode   int             `db:"keluaranperiode"`
	Datekeluaran      time.Time       `db:"datekeluaran"`
	Total_member      int             `db:"total_member"`
	Total_bet         decimal.Decimal `db:"total_bet"`
	Total_outstanding decimal.Decimal `db:"total_pairs"`
	Total_win         decimal.Decimal `db:"total_win"`
	Total_lose        decimal.Decimal `db:"total_lose"`
	Total_buangan     decimal.Decimal `db:"total_payout"`
	Total_reject      decimal.Decimal `db:"total_reject"`
}

// TrxkeluaranResultRow is one decided period (keluarantogel != ”) within a
// month range — the "Result" menu's past-draw-numbers list.
type TrxkeluaranResultRow struct {
	IDtrxkeluaran    int       `db:"idtrxkeluaran"`
	Keluaranperiode  int       `db:"keluaranperiode"`
	Datekeluaran     time.Time `db:"datekeluaran"`
	Keluarantogel    string    `db:"keluarantogel"`
	Aliascomppasaran string    `db:"aliascomppasaran"`
}

type TrxkeluaranRepository interface {
	FindAllRunning(ctx context.Context, idcomp string) ([]Trxkeluaranview, error)
	FindByID(ctx context.Context, idcomp, idcomppasaran string) (Trxkeluaran, error)
	FindByIDByNomorKeluaran(ctx context.Context, idcomp, idcomppasaran string) (Trxkeluaran, error)
	// IncrementTotals adds this checkout's contribution to the period's
	// running totals via a single "col = col + $delta" UPDATE rather than a
	// read-modify-write, so concurrent checkouts from many different players
	// against the same period can't stomp on each other — Postgres row-locks
	// the row for the statement's duration and serializes competing writers.
	IncrementTotals(ctx context.Context, idcomp string, idtrxkeluaran, totalMember int, totalBet, totalOutstanding, totalBuangan decimal.Decimal) error
	// FindResultsByMonth returns every period whose draw result has been
	// decided (keluarantogel != '') with datekeluaran in [start, end) for
	// idcomppasaran, newest first.
	FindResultsByMonth(ctx context.Context, idcomp, idcomppasaran string, start, end time.Time) ([]TrxkeluaranResultRow, error)
}
type TrxkeluaranService interface {
	All(ctx context.Context, idcomp string) ([]dto.TrxkeluaranData, error)
}
