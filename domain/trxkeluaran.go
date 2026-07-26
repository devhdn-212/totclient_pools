package domain

import (
	"context"
	"database/sql"
	"time"

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
	FindByID(ctx context.Context, idcomp, idcomppasaran string) (Trxkeluaran, error)
	FindByIDByNomorKeluaran(ctx context.Context, idcomp, idcomppasaran string) (Trxkeluaran, error)
	// RefreshTotals recomputes total_member/total_bet/total_pairs/
	// total_payout for idtrxkeluaran straight from the source rows in
	// tbl_trx_keluarantogel_member (COUNT DISTINCT username / SUM totalbet /
	// SUM totalpair / SUM totalbayar) and writes the absolute result — not an
	// increment. Self-healing by design: called repeatedly (see
	// service.FlushDirtyTotals) or out of order, it always converges on the
	// same correct totals with no drift, unlike the old per-checkout
	// increment which serialized every concurrent checkout for a period on
	// one contended row lock.
	RefreshTotals(ctx context.Context, idcomp string, idtrxkeluaran int) error
	// FindResultsByMonth returns every period whose draw result has been
	// decided (keluarantogel != '') with datekeluaran in [start, end) for
	// idcomppasaran, newest first.
	FindResultsByMonth(ctx context.Context, idcomp, idcomppasaran string, start, end time.Time) ([]TrxkeluaranResultRow, error)
}
