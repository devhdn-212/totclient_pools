package dto

type CurrData struct {
	ID      string `json:"curr_id"`
	Type    string `json:"curr_type"`
	Status  string `json:"curr_status"`
	Created string `json:"curr_created"`
	Update  string `json:"curr_updated"`
}
type CurrSave struct {
	Type      string `json:"type" validate:"required"`
	ID        string `json:"curr_id" validate:"required"`
	Type_curr string `json:"curr_type" validate:"required"`
	Status    string `json:"curr_status" validate:"required,max=1"`
}
