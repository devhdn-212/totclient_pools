package validasicustom

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

var uomRegex = regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)

func init() {
	Validate.RegisterValidation("uomname", func(fl validator.FieldLevel) bool {
		return uomRegex.MatchString(fl.Field().String())
	})
}
