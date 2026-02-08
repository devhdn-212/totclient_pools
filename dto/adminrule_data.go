package dto

type AdminruleData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
	Update  string `json:"updated"`
}
type AdminruleSelect struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type AdminruleSave struct {
	Type string `json:"type" validate:"required"`
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,max=30"`
	Rule string `json:"rule"`
}
