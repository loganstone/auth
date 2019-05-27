package validator

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// Validator ..
type Validator struct {
	validator *validator.Validate
}

// Validate ...
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err == nil {
		return err
	}

	he, ok := err.(*echo.HTTPError)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return he
}

// New ...
func New() *Validator {
	return &Validator{validator: validator.New()}
}
