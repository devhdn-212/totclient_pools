package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totagen_api/dto"
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
type TrxkeluaranmemberRepository interface {
	FindAll(ctx context.Context, idcomp string, idtrx int) ([]Trxkeluaranmember, error)
	FindByID(ctx context.Context, idcomp, idtrxdetail string, idtrx int) (Trxkeluaranmember, error)
	Save(ctx context.Context, trxkeluaranmember *Trxkeluaranmember, idcomp string) error
	Update(ctx context.Context, trxkeluaranmember *Trxkeluaranmember, idcomp string) error
}
type TrxkeluaranmemberService interface {
	All(ctx context.Context, idcomp string, idtrx int) ([]dto.TrxkeluaranmemberData, error)
	Save(ctx context.Context, req dto.TrxkeluaranmemberSave, client, idcomp string) error
}
