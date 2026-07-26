package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/devhdn-212/totclient_api/dto"
	"github.com/shopspring/decimal"
)

type Trxkeluaranmember struct {
	ID            string          `db:"idkeluaranmember"`
	IDtrxkeluaran int             `db:"idtrxkeluaran"`
	IDcomp        string          `db:"idcompany"`
	Username      string          `db:"username"`
	Totalbet      decimal.Decimal `db:"totalbet"`
	Totalbayar    decimal.Decimal `db:"totalbayar"`
	Totaldiscount decimal.Decimal `db:"totaldiscount"`
	Totalkei      decimal.Decimal `db:"totalkei"`
	Totalwin      decimal.Decimal `db:"totalwin"`
	Totalpair     int             `db:"totalpair"`
	Betround      int             `db:"betround"`
	Playerinvoice int             `db:"playerinvoice"`
	Status        string          `db:"status"`
	Created       string          `db:"createkeluaranmember"`
	CreatedAt     sql.NullTime    `db:"createdatekeluaranmember"`
	Update        string          `db:"updatekeluaranmember"`
	UpdateAt      sql.NullTime    `db:"updatedatekeluaranmember"`
}

// TrxkeluaranmemberPeriod is one aggregated row per idtrxkeluaran a player
// has ever transacted in for a pasaran — every playerinvoice's totalpair/
// totalbayar/totalwin summed. Used for the "Transaksi" period list (level
// 1); nothing here needs per-invoice granularity, so the aggregation
// happens in SQL instead of pulling every member row and grouping
// client-side. Keluarantogel empty means the draw for this period hasn't
// been decided yet (pending); non-empty means it has (complete).
type TrxkeluaranmemberPeriod struct {
	IDtrxkeluaran    int             `db:"idtrxkeluaran"`
	Datekeluaran     time.Time       `db:"datekeluaran"`
	Keluaranperiode  int             `db:"keluaranperiode"`
	Keluarantogel    string          `db:"keluarantogel"`
	Aliascomppasaran string          `db:"aliascomppasaran"`
	Codecomppasaran  string          `db:"codecomppasaran"`
	Totalpair        int             `db:"totalpair"`
	Totalbayar       decimal.Decimal `db:"totalbayar"`
	Totalwin         decimal.Decimal `db:"totalwin"`
}

type TrxkeluaranmemberRepository interface {
	FindAll(ctx context.Context, idcomp string, idtrx int) ([]Trxkeluaranmember, error)
	FindByID(ctx context.Context, idcomp, idtrxdetail string, idtrx int) (Trxkeluaranmember, error)
	// FindPeriodsByUsername returns one row per idtrxkeluaran the player has
	// ever transacted in for idcomppasaran, not just the currently-open one,
	// so a player's full "Transaksi" history survives period rollovers.
	FindPeriodsByUsername(ctx context.Context, idcomp, idcomppasaran, username string) ([]TrxkeluaranmemberPeriod, error)
	// ExistsByUsername reports whether username already has a member row for
	// idtrxkeluaran — checkout uses this (checked before it saves its own
	// member row) to decide whether this invoice is the player's first for
	// the period, so the period's total_member should count them once.
	ExistsByUsername(ctx context.Context, idcomp string, idtrxkeluaran int, username string) (bool, error)
	Save(ctx context.Context, trxkeluaranmember *Trxkeluaranmember, idcomp string) error
	Update(ctx context.Context, trxkeluaranmember *Trxkeluaranmember, idcomp string) error
}
type TrxkeluaranmemberService interface {
	All(ctx context.Context, idcomp string, idtrx int) ([]dto.TrxkeluaranmemberData, error)
	// Periods caches its result per idcomppasaran+username so the cache key
	// stays valid across period rollovers — checkout.go invalidates by the
	// same key, see trxkeluaranmemberByUserCacheKey.
	Periods(ctx context.Context, idcomp, idcomppasaran, username string) ([]dto.TrxkeluaranmemberPeriodData, error)
	Save(ctx context.Context, req dto.TrxkeluaranmemberSave, client, idcomp string) error
}
