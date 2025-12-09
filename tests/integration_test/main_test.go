package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDBPool *pgxpool.Pool //nolint:gochecknoglobals // Used by the integration test suite

//nolint:forbidigo // Allow fmt in test main
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Init postgres testcontainer
	pgTC, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Printf("failed to start container: %s\n", err)
		os.Exit(1)
	}

	conStr, _ := pgTC.ConnectionString(ctx, "sslmode=disable")

	// Open connection to run migrations
	gooseCon, err := sql.Open("pgx", conStr)
	if err != nil {
		fmt.Printf("failed to connect to database: %s\n", err)
		os.Exit(1)
	}

	if gErr := goose.Up(gooseCon, "../../internal/migrations"); gErr != nil {
		fmt.Printf("failed to run migrations: %s\n", gErr)
		os.Exit(1)
	}
	gooseCon.Close()

	// Open pgxpool for tests
	testDBPool, err = pgxpool.New(ctx, conStr)
	if err != nil {
		fmt.Printf("failed to create pgxpool: %s\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown
	testDBPool.Close()
	if err := pgTC.Terminate(ctx); err != nil {
		fmt.Printf("failed to terminate container: %s\n", err)
		os.Exit(1)
	}
	os.Exit(code)
}
