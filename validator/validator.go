package validator

import (
	"github.com/go-playground/validator"
)

// Validator ..
type Validator struct {
	validator *validator.Validate
}

// Validate ...
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// New ...
func New() *Validator {
	return &Validator{validator: validator.New()}
}
