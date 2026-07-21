package dto

type MemberinfoResponse struct {
	Agen   string `json:"agen" validate:"required"`
	Market string `json:"market" validate:"required"`
	Token  string `json:"token" validate:"required"`
}
