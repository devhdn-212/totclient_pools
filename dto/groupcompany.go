package dto

type GroupcompanyData struct {
	ID      string `json:"groupcompany_id"`
	Name    string `json:"groupcompany_name"`
	Status  string `json:"groupcompany_status"`
	Created string `json:"groupcompany_created"`
	Update  string `json:"groupcompany_updated"`
}
type GroupcompanySelect struct {
	ID   string `json:"groupcompany_id"`
	Name string `json:"groupcompany_name"`
}
type GroupcompanySave struct {
	Type   string `json:"type" validate:"required"`
	ID     string `json:"groupcompany_id" validate:"required"`
	Name   string `json:"groupcompany_name" validate:"required,max=100"`
	Status string `json:"groupcompany_status" validate:"required,max=1"`
}
