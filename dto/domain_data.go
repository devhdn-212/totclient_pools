package dto

type DomainData struct {
	ID      string `json:"domain_id"`
	Type    string `json:"domain_type"`
	Name    string `json:"domain_name"`
	Status  string `json:"domain_status"`
	Created string `json:"domain_created"`
	Update  string `json:"domain_updated"`
}
type DomainSave struct {
	Type       string `json:"type" validate:"required"`
	ID         string `json:"domain_id" `
	Typedomain string `json:"domain_type" validate:"required"`
	Name       string `json:"domain_name" validate:"required"`
	Status     string `json:"domain_status" validate:"required,max=1"`
}
