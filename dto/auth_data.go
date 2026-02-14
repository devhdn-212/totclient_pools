package dto

type AuthRequest struct {
	Username  string `json:"username" validate:"required,min=4,max=30,lowercase,alphanum"`
	Password  string `json:"password" validate:"required"`
	Ipaddress string `json:"ipaddress" validate:"required"`
}
type AuthPage struct {
	Page string `json:"page" validate:"required"`
}
type AuthClientRedis struct {
	Username string `json:"username"`
	IDrule   string `json:"idule"`
	Rule     string `json:"rule"`
}
type AuthResponse struct {
	Token string `json:"token"`
}
