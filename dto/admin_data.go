package dto

type AdminData struct {
	ID        string `json:"admin_id"`
	Username  string `json:"admin_username"`
	Idrule    string `json:"admin_rule"`
	Name      string `json:"admin_name"`
	Lastlogin string `json:"admin_lastlogin"`
	Joindate  string `json:"admin_joindate"`
	Ipaddress string `json:"admin_ipaddress"`
	Status    string `json:"admin_status"`
	Created   string `json:"admin_created"`
	Update    string `json:"admin_updated"`
}
type AdminLogin struct {
	Username string `json:"username" validate:"required"`
	Pass     string `json:"password" validate:"required"`
}

type AdminSave struct {
	Type     string `json:"type" validate:"required"`
	ID       string `json:"admin_id" validate:"required"`
	IDrule   string `json:"admin_idrule" validate:"required"`
	Username string `json:"admin_username" validate:"required"`
	Pass     string `json:"admin_password" `
	Name     string `json:"admin_name" validate:"required"`
	Status   string `json:"admin_status" validate:"required,max=1"`
}
