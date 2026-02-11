package dto

type ClientruleData struct {
	ID      string `json:"clientrule_id"`
	Name    string `json:"clientrule_name"`
	Rule    string `json:"clientrule_rule"`
	Created string `json:"clientrule_create"`
	Update  string `json:"clientrule_update"`
}
type ClientruleSelect struct {
	ID   string `json:"clientrule_id"`
	Name string `json:"clientrule_name"`
}
type ClientruleSave struct {
	Type string `json:"type" validate:"required"`
	ID   string `json:"clientrule_id" validate:"required"`
	Name string `json:"clientrule_name" validate:"required,max=30"`
	Rule string `json:"clientrule_rule"`
}
