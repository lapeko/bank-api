package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

func CurrencyValidator(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(utils.Currency); ok {
		return utils.IsCurrency(currency)
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
