package dto

import "github.com/shopspring/decimal"

type CompanyData struct {
	ID        string          `json:"company_id"`
	IDcurrdef string          `json:"company_idcurr"`
	Name      string          `json:"company_name"`
	Endjoin   string          `json:"company_endjoin"`
	Amount    decimal.Decimal `json:"company_amount"`
	Status    string          `json:"company_status"`
	Created   string          `json:"company_created"`
	Update    string          `json:"company_updated"`
}
type CompanySave struct {
	Type   string `json:"type" validate:"required"`
	ID     string `json:"company_id" validate:"required"`
	IDcurr string `json:"company_idcurr" validate:"required"`
	Name   string `json:"company_name" validate:"required,max=50"`
	Status string `json:"company_status" validate:"required,max=1"`
}
