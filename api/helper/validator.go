package helper

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateSlug(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(string); ok {
		return regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`).MatchString(t)
	}
	return false
}
