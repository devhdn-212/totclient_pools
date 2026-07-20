package dto

import "github.com/shopspring/decimal"

type TrxkeluaranData struct {
	ID                int             `json:"trxkeluaran_id"`
	IDcompasaran      string          `json:"trxkeluaran_idcomppasaran"`
	IDcomp            string          `json:"trxkeluaran_idcompany"`
	Nmpasaran         string          `json:"trxkeluaran_nmpasaran"`
	Yearmonth         string          `json:"trxkeluaran_yearmonth"`
	Periode           int             `json:"trxkeluaran_keluaranperiode"`
	Datekeluaran      string          `json:"trxkeluaran_datekeluaran"`
	Keluarantogel     string          `json:"trxkeluaran_keluarantogel"`
	Prize2            string          `json:"trxkeluaran_prize2"`
	Prize3            string          `json:"trxkeluaran_prize3"`
	Total_member      int             `json:"trxkeluaran_total_member"`
	Total_bet         decimal.Decimal `json:"trxkeluaran_total_bet"`
	Total_outstanding decimal.Decimal `json:"trxkeluaran_total_outstanding"`
	Total_win         decimal.Decimal `json:"trxkeluaran_total_win"`
	Total_lose        decimal.Decimal `json:"trxkeluaran_total_lose"`
	Total_buangan     decimal.Decimal `json:"trxkeluaran_total_buangan"`
	Total_reject      decimal.Decimal `json:"trxkeluaran_total_reject"`
	Winlose           decimal.Decimal `json:"trxkeluaran_winlose"`
	Revisi            decimal.Decimal `json:"trxkeluaran_revisi"`
	Noterevisi        decimal.Decimal `json:"trxkeluaran_noterevisi"`
	Created           string          `json:"trxkeluaran_created"`
	Update            string          `json:"trxkeluaran_updated"`
}

type TrxkeluaranSave struct {
	Type         string `json:"type" validate:"required"`
	IDcompasaran string `json:"trxkeluaran_idcomppasaran" validate:"required"`
}
