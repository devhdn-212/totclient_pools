package dto

type AdminruleData struct {
	ID      string `json:"adminrule_id"`
	Name    string `json:"adminrule_name"`
	Rule    string `json:"adminrule_rule"`
	Created string `json:"adminrule_create"`
	Update  string `json:"adminrule_update"`
}
type AdminruleSelect struct {
	ID   string `json:"adminrule_id"`
	Name string `json:"adminrule_name"`
}
type AdminruleSave struct {
	Type string `json:"type" validate:"required"`
	ID   string `json:"adminrule_id" validate:"required"`
	Name string `json:"adminrule_name" validate:"required,max=30"`
	Rule string `json:"adminrule_rule"`
}
