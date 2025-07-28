package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
)

func RegisterValidators() {
	once.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation("currency", currencyValidator)
			v.RegisterValidation("fullname", fullNameValidator)
			v.RegisterValidation("password", passwordValidator)
		}
	})
}

func currencyValidator(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(utils.Currency); ok {
		return utils.IsCurrency(currency)
	}
	return false
}

func fullNameValidator(fl validator.FieldLevel) bool {
	if fullName, ok := fl.Field().Interface().(string); ok {
		return len(fullName) >= 4
	}
	return false
}

func passwordValidator(fl validator.FieldLevel) bool {
	if password, ok := fl.Field().Interface().(string); ok {
		return len(password) >= 6
	}
	return false
}
