package dto

import "github.com/shopspring/decimal"

type TrxkeluarandetailData struct {
	ID                   string          `json:"trxkeluarandetail_id"`
	IDtrxkeluaran        int             `json:"trxkeluarandetail_idtrxkeluaran"`
	IDcomp               string          `json:"trxkeluarandetail_idcompany"`
	Datekeluarandetail   string          `json:"trxkeluarandetail_datetimedetail"`
	Ipaddress            string          `json:"trxkeluarandetail_ipaddress"`
	Username             string          `json:"trxkeluarandetail_username"`
	Typegame             string          `json:"trxkeluarandetail_typegame"`
	Nomortogel           string          `json:"trxkeluarandetail_nomortogel"`
	Posisitogel          string          `json:"trxkeluarandetail_posisitogel"`
	Bet                  int             `json:"trxkeluarandetail_bet"`
	Diskon               decimal.Decimal `json:"trxkeluarandetail_diskon"`
	Kei                  decimal.Decimal `json:"trxkeluarandetail_kei"`
	Win                  decimal.Decimal `json:"trxkeluarandetail_win"`
	Winhasil             decimal.Decimal `json:"trxkeluarandetail_winhasil"`
	Winrev               decimal.Decimal `json:"trxkeluarandetail_winrev"`
	Betround             int             `json:"trxkeluarandetail_betround"`
	Browser              string          `json:"trxkeluarandetail_browsertogel"`
	Devicetogel          string          `json:"trxkeluarandetail_devicetogel"`
	Playerinvoice        int             `json:"trxkeluarandetail_playerinvoice"`
	Statuskeluarandetail string          `json:"trxkeluarandetail_statuskeluarandetail"`
	Senddata             string          `json:"trxkeluarandetail_senddata"`
	Updatedata           string          `json:"trxkeluarandetail_updatedata"`
	Created              string          `json:"trxkeluarandetail_created"`
	Update               string          `json:"trxkeluarandetail_updated"`
}
type TrxkeluarandetailAll struct {
	IDtrxkeluaran int `json:"trxkeluarandetail_idtrxkeluaran" validate:"required"`
}
type TrxkeluarandetailSave struct {
	Type          string `json:"type" validate:"required"`
	ID            string `json:"trxkeluarandetail_id" `
	IDtrxkeluaran int    `json:"trxkeluarandetail_idtrxkeluaran" validate:"required"`
}
