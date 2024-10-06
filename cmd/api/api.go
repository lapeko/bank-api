package main

import (
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/config"
)

func main() {
	conf := config.NewApiConfig()
	conf.Parse()
	a := api.New()
	a.ConnectStore(conf)
	a.SetUpRoutes()
	a.Start()
}
