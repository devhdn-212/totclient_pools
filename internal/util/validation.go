package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := map[string]string{}
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			res[v.StructField()] = TransalateTaq(v)
		}
	}
	return res
}
func TransalateTaq(fd validator.FieldError) string {
	switch fd.Tag() {
	case "required":
		return fmt.Sprintf("Field %s is required", fd.Field())
	case "lowercase":
		return fmt.Sprintf("Field %s must be lowercase", fd.Field())
	case "alphanum":
		return fmt.Sprintf("Field %s may contain only letters and numbers", fd.Field())
	}
	return "validasi failed"
}
