package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totclient_api/dto"
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
	Cancelbet            decimal.Decimal `db:"cancelbet"`
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
	FindByUsername(ctx context.Context, idcomp string, idtrx int, username string) ([]Trxkeluarandetail, error)
	FindByID(ctx context.Context, idcomp, idtrxdetail string, idtrx int) (Trxkeluarandetail, error)
	Save(ctx context.Context, trxkeluarandetail *Trxkeluarandetail, idcomp string) error
}
type TrxkeluarandetailService interface {
	All(ctx context.Context, idcomp string, idtrx int) ([]dto.TrxkeluarandetailData, error)
	AllByUsername(ctx context.Context, idcomp string, idtrx int, username string) ([]dto.TrxkeluarandetailData, error)
	Save(ctx context.Context, req dto.TrxkeluarandetailSave, client, idcomp string) error
}
