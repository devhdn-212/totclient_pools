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
type TrxkeluaranmemberAll struct {
	IDtrxkeluaran int `json:"trxkeluarandetail_idtrxkeluaran" validate:"required"`
}
type TrxkeluaranmemberSave struct {
	Type          string `json:"type" validate:"required"`
	ID            string `json:"trxkeluaramember_id" `
	IDtrxkeluaran int    `json:"trxkeluaranmember_idtrxkeluaran" validate:"required"`
	Playerinvoice int    `json:"trxkeluaranmember_playerinvoice" validate:"required"`
}
