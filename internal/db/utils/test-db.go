package utils

import (
	"context"
	"fmt"
	logOrigin "log"
	"os"
	"os/exec"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	exposedPort = "5432/tcp"
	dbUser      = "postgres"
	dbPass      = "1234"
	dbName      = "bank_test"
)

var errLog = logOrigin.New(os.Stderr, "[test-utils] testcontainer initialization failure. Error: ", logOrigin.Lshortfile)

func SetupTestDb(ctx context.Context) (string, func() error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:alpine",
			ExposedPorts: []string{exposedPort},
			Env: map[string]string{
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPass,
				"POSTGRES_DB":       dbName,
			},
			WaitingFor: wait.ForListeningPort(exposedPort),
		},
		Started: true,
	})

	if err != nil {
		errLog.Fatalf("%q", err)
	}

	terminate := func() error {
		return container.Terminate(ctx)
	}

	host, err := container.Host(ctx)
	if err != nil {
		if e := terminate(); e != nil {
			errLog.Printf("%q", e)
		}
		errLog.Fatalf("%q", err)
	}

	port, err := container.MappedPort(ctx, exposedPort)
	if err != nil {
		if e := terminate(); e != nil {
			errLog.Printf("%q", e)
		}
		errLog.Fatalf("%q", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, host, port.Port(), dbName)

	cmd := exec.Command("migrate", "-source", "file://../migration", "-database", dsn, "-verbose", "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[migrate] %s\n", output)
		errLog.Fatalf("migrate command failed: %v", err)
	}

	return dsn, terminate
}
