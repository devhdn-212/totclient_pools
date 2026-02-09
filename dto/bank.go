package dto

type BankData struct {
	ID      string `json:"bank_id"`
	Type    string `json:"bank_type"`
	Name    string `json:"bank_name"`
	Status  string `json:"bank_status"`
	Created string `json:"bank_created"`
	Update  string `json:"bank_updated"`
}
type BankSelect struct {
	ID   string `json:"bank_id"`
	Name string `json:"bank_name"`
}
type BankSave struct {
	Type      string `json:"type" validate:"required"`
	ID        string `json:"bank_id" validate:"required"`
	Type_bank string `json:"bank_type" validate:"required"`
	Name      string `json:"bank_name" validate:"required"`
	Status    string `json:"bank_status" validate:"required,max=1"`
}
