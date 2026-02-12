package dto

type CompanyadminData struct {
	ID        string `json:"compadmin_id"`
	IDcomp    string `json:"compadmin_idccomp"`
	Rule      string `json:"compadmin_rule"`
	Username  string `json:"compadmin_username"`
	Name      string `json:"compadmin_name"`
	Lastlogin string `json:"compadmin_lastlogin"`
	Ipaddress string `json:"compadmin_ipaddress"`
	Status    string `json:"compadmin_status"`
	Created   string `json:"compadmin_created"`
	Update    string `json:"compadmin_updated"`
}
type CompanyadminAll struct {
	IDcompany string `json:"compadmin_idcompany" validate:"required"`
}
type CompanyadminSave struct {
	Type      string `json:"type" validate:"required"`
	ID        string `json:"compadmin_id"`
	IDcompany string `json:"compadmin_idcompany" validate:"required"`
	IDrule    string `json:"compadmin_idrule" validate:"required"`
	Username  string `json:"compadmin_username" validate:"required"`
	Pass      string `json:"compadmin_password" `
	Name      string `json:"compadmin_name" validate:"required"`
	Status    string `json:"compadmin_status" validate:"required,max=1"`
}
