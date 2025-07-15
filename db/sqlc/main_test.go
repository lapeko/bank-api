package db

import (
	"context"
	logOrigin "log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/db/utils"
)

var (
	testStore Store
	log       = logOrigin.New(os.Stderr, "[db] ", logOrigin.Lshortfile)
	ctx       context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	dsn, terminate := utils.SetupTestDb(ctx)
	dbConnection, err := pgxpool.New(ctx, dsn)

	if err != nil {
		if e := terminate(); e != nil {
			log.Printf("container termination failure. Error: %q", e)
		}
		log.Fatalf("connection failure. Error: %q", err)
	}

	testStore = NewStore(dbConnection)
	code := m.Run()

	dbConnection.Close()

	if err := terminate(); err != nil {
		log.Printf("container termination failure. Error: %q", err)
		code = 1
	}

	os.Exit(code)
}
