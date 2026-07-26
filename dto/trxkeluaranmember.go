package dto

import "github.com/shopspring/decimal"

type TrxkeluaranmemberData struct {
	ID            string          `json:"trxkeluaranmember_id"`
	IDtrxkeluaran int             `json:"trxkeluaranmember_idtrxkeluaran"`
	IDcomp        string          `json:"trxkeluaranmember_idcompany"`
	Username      string          `json:"trxkeluaranmember_username"`
	Totalbet      decimal.Decimal `json:"trxkeluaranmember_totalbet"`
	Totalbayar    decimal.Decimal `json:"trxkeluaranmember_totalbayar"`
	Totaldiscount decimal.Decimal `json:"trxkeluaranmember_totaldiscount"`
	Totalkei      decimal.Decimal `json:"trxkeluaranmember_totalkei"`
	Totalwin      decimal.Decimal `json:"trxkeluaranmember_totalwin"`
	Totalpair     int             `json:"trxkeluaranmember_totalpair"`
	Betround      int             `json:"trxkeluaranmember_betround"`
	Playerinvoice int             `json:"trxkeluaranmember_playerinvoice"`
	Status        string          `json:"trxkeluaranmember_status"`
	Created       string          `json:"trxkeluaranmember_created"`
	Update        string          `json:"trxkeluaranmember_updated"`
}

// TrxkeluaranmemberPeriodData is one aggregated row per idtrxkeluaran a
// player has ever transacted in for a pasaran — the "Transaksi" period
// list (level 1). Status is "pending" (draw not decided yet), "WIN", or
// "LOSE" (draw decided).
type TrxkeluaranmemberPeriodData struct {
	IDtrxkeluaran    int             `json:"idtrxkeluaran"`
	Datekeluaran     string          `json:"datekeluaran"`
	Keluaranperiode  int             `json:"keluaranperiode"`
	Aliascomppasaran string          `json:"aliascomppasaran"`
	Codecomppasaran  string          `json:"codecomppasaran"`
	Totalpair        int             `json:"totalpair"`
	Totalbayar       decimal.Decimal `json:"totalbayar"`
	Totalwin         decimal.Decimal `json:"totalwin"`
	Status           string          `json:"status"`
}

type TrxkeluaranmemberAll struct {
	IDtrxkeluaran int `json:"trxkeluarandetail_idtrxkeluaran" validate:"required"`
}
type TrxkeluaranmemberSave struct {
	Type          string `json:"type" validate:"required"`
	ID            string `json:"trxkeluaramember_id" `
	IDtrxkeluaran int    `json:"trxkeluaranmember_idtrxkeluaran" validate:"required"`
	Playerinvoice int    `json:"trxkeluaranmember_playerinvoice" validate:"required"`
}
