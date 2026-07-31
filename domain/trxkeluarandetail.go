package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/devhdn-212/totclient_pools/dto"
	"github.com/shopspring/decimal"
)

type Trxkeluarandetail struct {
	ID                   string          `db:"idtrxkeluarandetail"`
	IDtrxkeluaran        int             `db:"idtrxkeluaran"`
	IDcomp               string          `db:"idcompany"`
	Datekeluarandetail   sql.NullTime    `db:"datetimedetail"`
	Ipaddress            string          `db:"ipaddress"`
	Username             string          `db:"username"`
	Typegame             string          `db:"typegame"`
	Nomortogel           string          `db:"nomortogel"`
	Posisitogel          string          `db:"posisitogel"`
	Bet                  int             `db:"bet"`
	Diskon               decimal.Decimal `db:"diskon"`
	Win                  decimal.Decimal `db:"win"`
	Winhasil             decimal.Decimal `db:"winhasil"`
	Kei                  decimal.Decimal `db:"kei"`
	Browsertogel         string          `db:"browsertogel"`
	Devicetogel          string          `db:"devicetogel"`
	Statuskeluarandetail string          `db:"statuskeluarandetail"`
	Betround             int             `db:"betround"`
	Winrev               decimal.Decimal `db:"winrev"`
	Playerinvoice        int             `db:"playerinvoice"`
	Senddata             string          `db:"senddata"`
	Senddatacreatedate   sql.NullTime    `db:"senddatacreatedate"`
	Updatedata           string          `db:"updatedata"`
	Updatedatacreatedate sql.NullTime    `db:"updatedatacreatedate"`
	Created              string          `db:"create_by"`
	CreatedAt            sql.NullTime    `db:"create_at"`
	Update               string          `db:"update_by"`
	UpdateAt             sql.NullTime    `db:"update_at"`
}
type TrxkeluarandetailRepository interface {
	FindAll(ctx context.Context, idcomp string, idtrx int) ([]Trxkeluarandetail, error)
	// FindByUsername matches against every idtrxkeluaran in idtrx — the
	// caller sources that list from TrxkeluaranmemberRepository.FindAllByUsername
	// (every period the player was active in for a pasaran), so this covers
	// the player's full "Transaksi" history across period rollovers.
	FindByUsername(ctx context.Context, idcomp string, idtrx []int, username string) ([]Trxkeluarandetail, error)
	FindByID(ctx context.Context, idcomp, idtrxdetail string, idtrx int) (Trxkeluarandetail, error)
	Save(ctx context.Context, trxkeluarandetail *Trxkeluarandetail, idcomp string) error
	// SumBet totals up bet across every already-persisted row matching
	// (idtrx, typegame, nomortogel), optionally narrowed to one username —
	// used to reseed a limittotal/limitglobal Redis counter from the
	// permanent ledger when Redis has no memory of it (fresh key, e.g. after
	// a Redis restart). username == "" sums across every player (limitglobal);
	// a non-empty username scopes it to just that player (limittotal).
	// datekeluaran is the draw date for idtrx — the table is partitioned by
	// month on datetimedetail, so passing it lets Postgres prune straight to
	// the relevant partition(s) instead of scanning every partition ever
	// created.
	SumBet(ctx context.Context, idcomp string, idtrx int, datekeluaran time.Time, username, typegame, nomortogel string) (int64, error)
}
type TrxkeluarandetailService interface {
	All(ctx context.Context, idcomp string, idtrx int) ([]dto.TrxkeluarandetailData, error)
	// AllByUsername is not cached — idtrx is normally a single period the
	// client just asked to open, so it varies call-to-call for the same
	// idcomppasaran+username. See the implementation for why.
	AllByUsername(ctx context.Context, idcomp, idcomppasaran string, idtrx []int, username string) ([]dto.TrxkeluarandetailData, error)
	Save(ctx context.Context, req dto.TrxkeluarandetailSave, client, idcomp string) error
}
