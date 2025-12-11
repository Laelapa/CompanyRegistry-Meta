package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/app"
	"github.com/Laelapa/CompanyRegistry/internal/config"
	"github.com/Laelapa/CompanyRegistry/internal/repository"
	"github.com/Laelapa/CompanyRegistry/internal/repository/adapters"
	"github.com/Laelapa/CompanyRegistry/internal/service"
	"github.com/Laelapa/CompanyRegistry/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDBPool *pgxpool.Pool //nolint:gochecknoglobals // Used by the integration test suite

func TestMain(m *testing.M) {
	// Delegate to a helper function so we can use defer for cleanup
	code := tRun(m)
	os.Exit(code)
}

//nolint:forbidigo // Allow fmt in test setup
func tRun(m *testing.M) int {
	ctx := context.Background()

	// This is used for compatibility issues with specific OSes,
	// Try to run without it first, only use if necessary
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

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
		return 1
	}

	// Defer termination immediately so it runs even if setup fails later
	defer func() {
		if tcerr := pgTC.Terminate(ctx); tcerr != nil {
			fmt.Printf("failed to terminate container: %s\n", tcerr)
		} else {
			fmt.Println("Testcontainer terminated successfully")
		}
	}()

	conStr, _ := pgTC.ConnectionString(ctx, "sslmode=disable")

	// Open connection to run migrations
	gooseCon, err := sql.Open("pgx", conStr)
	if err != nil {
		fmt.Printf("failed to connect to database: %s\n", err)
		return 1
	}
	defer gooseCon.Close()

	if gErr := goose.Up(gooseCon, "../../internal/migrations"); gErr != nil {
		fmt.Printf("failed to run migrations: %s\n", gErr)
		return 1
	}

	// Open pgxpool for tests
	testDBPool, err = pgxpool.New(ctx, conStr)
	if err != nil {
		fmt.Printf("failed to create pgxpool: %s\n", err)
		return 1
	}
	defer testDBPool.Close()

	// Run tests
	return m.Run()
}

func setupApp(t *testing.T) *app.App {
	t.Helper()

	// TODO: develop test logger
	logger, _ := logging.NewLogger(config.LoggingConfig{LoggerSetup: "prod"})
	tokenAuth := tokenauthority.New(
		&config.AuthConfig{
			JwtSecret:   "test-secret-key",
			JwtIssuer:   "test-issuer",
			JwtLifetime: 1 * time.Hour,
		},
	)
	queries := repository.New(testDBPool)

	svc := &service.Service{
		User: service.NewUserService(
			adapters.NewPGUserRepoAdapter(queries),
			tokenAuth,
			logger,
			nil,
			"doesn't-matter",
		),
		Company: service.NewCompanyService(
			adapters.NewPGCompanyRepoAdapter(queries),
			tokenAuth,
			logger,
			nil,
			"doesn't-matter",
		),
	}
	srvCfg := &config.ServerConfig{
		Port:            "8080",
		ShutdownTimeout: 5 * time.Second,
	}

	return app.New(
		srvCfg,
		logger,
		svc,
		tokenAuth,
		nil,
	)
}
