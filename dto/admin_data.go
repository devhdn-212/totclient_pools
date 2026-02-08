package dto

type AdminData struct {
	Username  string `json:"admin_username"`
	Idadmin   string `json:"admin_role"`
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
	Type      string `json:"type" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Pass      string `json:"password" `
	Idadmin   string `json:"idadmin" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Ipaddress string `json:"ipaddress"`
	Status    string `json:"status" validate:"required,max=1"`
}
