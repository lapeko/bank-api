package utils

import "github.com/go-playground/validator/v10"

func CurrencyValidator(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(Currency); ok {
		return IsCurrency(currency)
	}
	return false
}

func FullNameValidator(fl validator.FieldLevel) bool {
	if fullName, ok := fl.Field().Interface().(string); ok {
		return len(fullName) >= 4
	}
	return false
}

func PasswordValidator(fl validator.FieldLevel) bool {
	if password, ok := fl.Field().Interface().(string); ok {
		return len(password) >= 6
	}
	return false
}
