package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/devhdn-212/totagen_api/dto"
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
	Total_outstanding decimal.Decimal `db:"total_outstanding"`
	Total_win         decimal.Decimal `db:"total_win"`
	Total_lose        decimal.Decimal `db:"total_lose"`
	Total_buangan     decimal.Decimal `db:"total_buangan"`
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
	Total_outstanding decimal.Decimal `db:"total_outstanding"`
	Total_win         decimal.Decimal `db:"total_win"`
	Total_lose        decimal.Decimal `db:"total_lose"`
	Total_buangan     decimal.Decimal `db:"total_buangan"`
	Total_reject      decimal.Decimal `db:"total_reject"`
}

type TrxkeluaranRepository interface {
	FindAllRunning(ctx context.Context, idcomp string) ([]Trxkeluaranview, error)
	FindByID(ctx context.Context, idcomp, idcomppasaran string, idtrx int) (Trxkeluaran, error)
	FindByIDByNomorKeluaran(ctx context.Context, idcomp, idcomppasaran string) (Trxkeluaran, error)
	Save(ctx context.Context, trxkeluaran *Trxkeluaran, idcomp string) error
	Update(ctx context.Context, trxkeluaran *Trxkeluaran, idcomp string) error
}
type TrxkeluaranService interface {
	All(ctx context.Context, idcomp string) ([]dto.TrxkeluaranData, error)
	Save(ctx context.Context, req dto.TrxkeluaranSave, client, idcomp string) error
}
