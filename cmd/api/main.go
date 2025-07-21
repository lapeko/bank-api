package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api"
)

func main() {
	postgresUrl := os.Getenv("POSTGRES_URL")
	portStr := os.Getenv("APP_PORT")
	if portStr == "" {
		portStr = "3000"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("unable to parse int port")
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, postgresUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}
	api := api.New(conn)
	log.Fatalln(api.Start(port))
}
