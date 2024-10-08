package main

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/config"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/validators"
)

func main() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.Currency)
	}
	conf := config.NewApiConfig()
	conf.Parse()
	a := api.GetApi()
	a.ConnectStore(conf)
	a.SetUpRoutes()
	a.Start()
}
