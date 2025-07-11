package db

import (
	"context"
	logOrigin "log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/db/utils"
)

const (
	dbPort = 5433
	dbPass = "1234"
	dbName = "bank_test"
	dbUser = "postgres"
)

var (
	sqlcQueries *Queries
	log         = logOrigin.New(os.Stderr, "[db] ", logOrigin.Lshortfile)
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn, terminate := utils.SetupTestDb(ctx)
	dbConnection, err := pgx.Connect(ctx, dsn)

	if err != nil {
		if e := terminate(); e != nil {
			log.Printf("container termination failure. Error: %q", e.Error())
		}
		log.Fatalf("connection failure. Error: %q", err.Error())
	}

	sqlcQueries = New(dbConnection)
	code := m.Run()

	if err := dbConnection.Close(ctx); err != nil {
		log.Fatalf("connection close failure. Error: %q", err.Error())
	}

	os.Exit(code)
}
