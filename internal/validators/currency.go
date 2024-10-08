package validators

import (
	"github.com/go-playground/validator/v10"
)

const (
	EUR = "EUR"
	USD = "USD"
	PLN = "PLN"
)

var Currency validator.Func = func(fl validator.FieldLevel) bool {
	curr, ok := fl.Field().Interface().(string)
	if ok {
		switch curr {
		case EUR, USD, PLN:
			return true
		default:
			return false
		}
	}
	return false
}
