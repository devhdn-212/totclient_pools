package dto

type CustomerData struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_datetime"`
	UpdateAt  string `json:"updated_datetime"`
}

type CreateCustomerRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}
type UpdateCustomerRequest struct {
	ID   string `json:"-"`
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}
