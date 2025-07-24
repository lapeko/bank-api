package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/config"
)

func main() {
	cfg := config.Get()
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, cfg.PostgresUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}
	api := api.New(conn)
	log.Fatalln(api.Start(cfg.AppPort))
}
