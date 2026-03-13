package dto

type UomData struct {
	ID      string `json:"uom_id"`
	Name    string `json:"uom_name"`
	Status  string `json:"uom_status"`
	Created string `json:"uom_created"`
	Update  string `json:"uom_updated"`
}
type UomSelect struct {
	ID   string `json:"uom_id"`
	Name string `json:"uom_name"`
}
type UomSave struct {
	Type   string `json:"type" validate:"required"`
	ID     string `json:"uom_id" validate:"required,min=1,max=10"`
	Name   string `json:"uom_name" validate:"required,min=2,max=100,uomname"`
	Status string `json:"uom_status" validate:"required,max=1"`
}
