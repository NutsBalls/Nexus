package services

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator will be used to validate structs with the `validate` tag
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
