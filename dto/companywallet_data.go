package dto

import "github.com/shopspring/decimal"

type CompanywalletData struct {
	ID      string          `json:"companywallet_id"`
	IDcomp  string          `json:"companywallet_idccomp"`
	IDcurr  string          `json:"companywallet_idcurr"`
	Amount  decimal.Decimal `json:"companywallet_amount"`
	Status  string          `json:"companywallet_status"`
	Created string          `json:"companywallet_created"`
	Update  string          `json:"companywallet_updated"`
}

type CompanywalletSave struct {
	Type   string `json:"type" validate:"required"`
	ID     string `json:"companywallet_id" `
	IDcomp string `json:"companywallet_idcompany" validate:"required"`
	IDcurr string `json:"companywallet_idcurr" validate:"required"`
	Status string `json:"companywallet_status" validate:"required,max=1"`
}
type CompanywalletReqIndex struct {
	IDcomp string `json:"companywallet_idcompany" validate:"required"`
}
