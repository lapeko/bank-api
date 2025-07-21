package utils

import (
	"github.com/go-playground/validator/v10"
)

type Currency string

const (
	CurrencyUSD  Currency = "USD"
	CurrencyEURO Currency = "EURO"
	CurrencyPLN  Currency = "PLN"
)

func IsCurrency(currency Currency) bool {
	switch currency {
	case CurrencyUSD, CurrencyEURO, CurrencyPLN:
		return true
	default:
		return false
	}
}

func CurrencyValidator(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(Currency); ok {
		return IsCurrency(currency)
	}
	return false
}
