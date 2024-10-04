package main

import "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api"

const (
	postgresAddr = "postgres://root:1234@127.0.0.1:5432/simple_bank?sslmode=disable"
	apiAddr      = ":8080"
)

func main() {
	a := api.New(postgresAddr)
	a.Start(apiAddr)
}
