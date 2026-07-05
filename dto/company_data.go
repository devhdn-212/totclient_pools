package dto

import "github.com/shopspring/decimal"

type CompanyData struct {
	ID          string          `json:"company_id"`
	IDgroupcomp string          `json:"company_idgroupcomp"`
	Nmgroupcomp string          `json:"company_nmgroupcomp"`
	IDcurrdef   string          `json:"company_idcurr"`
	Name        string          `json:"company_name"`
	Endjoin     string          `json:"company_endjoin"`
	Amount      decimal.Decimal `json:"company_amount"`
	TelegramID  string          `json:"company_telegramid"`
	URLapitoto  string          `json:"company_urlapitoto"`
	URLapislot  string          `json:"company_urlapislot"`
	Status      string          `json:"company_status"`
	Activetoto  string          `json:"company_activetoto"`
	Activeslot  string          `json:"company_activeslot"`
	Created     string          `json:"company_created"`
	Update      string          `json:"company_updated"`
}
type CompanySave struct {
	Type        string `json:"type" validate:"required"`
	ID          string `json:"company_id" validate:"required"`
	IDgroupcomp string `json:"company_idgroupcomp" validate:"required"`
	IDcurr      string `json:"company_idcurr" validate:"required"`
	Name        string `json:"company_name" validate:"required,max=50"`
	TelegramID  string `json:"company_telegramid" validate:"max=50"`
	URLapitoto  string `json:"company_urlapitoto" validate:"max=150"`
	URLapislot  string `json:"company_urlapislot" validate:"max=150"`
	Activetoto  string `json:"company_activetoto" `
	Activeslot  string `json:"company_activeslot" `
	Status      string `json:"company_status" validate:"required,max=1"`
}
