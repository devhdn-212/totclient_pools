package domain

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

// Trxkeluarandetailreject records a bet that failed limittotal or
// limitglobal at checkout time — kept for audit (so a rejected bet is never
// silently dropped) rather than mixed into Trxkeluarandetail with a status
// flag, matching the separate tbl_trx_keluarantogel_detail_reject table.
type Trxkeluarandetailreject struct {
	ID                 string          `db:"idtrxkeluarandetailreject"`
	IDtrxkeluaran      int             `db:"idtrxkeluaran"`
	IDcomp             string          `db:"idcompany"`
	Datekeluarandetail sql.NullTime    `db:"datetimedetail"`
	Ipaddress          string          `db:"ipaddress"`
	Username           string          `db:"username"`
	Typegame           string          `db:"typegame"`
	Nomortogel         string          `db:"nomortogel"`
	Posisitogel        string          `db:"posisitogel"`
	Bet                int             `db:"bet"`
	Playerinvoice      int             `db:"playerinvoice"`
	Reason             string          `db:"reason"`
	Sisalimit          decimal.Decimal `db:"sisalimit"`
	Browsertogel       string          `db:"browsertogel"`
	Devicetogel        string          `db:"devicetogel"`
	Created            string          `db:"create_by"`
	CreatedAt          sql.NullTime    `db:"create_at"`
}

type TrxkeluarandetailrejectRepository interface {
	Save(ctx context.Context, reject *Trxkeluarandetailreject, idcomp string) error
}
