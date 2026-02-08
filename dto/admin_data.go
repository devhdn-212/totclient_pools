package dto

type AdminData struct {
	Username  string `json:"username"`
	Idadmin   string `json:"idadmin"`
	Name      string `json:"name"`
	Lastlogin string `json:"lastlogin"`
	Joindate  string `json:"joindate"`
	Ipaddress string `json:"ipaddress"`
	Status    string `json:"status"`
	Created   string `json:"created"`
	Update    string `json:"updated"`
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
