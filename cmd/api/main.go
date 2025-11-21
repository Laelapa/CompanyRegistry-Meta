package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Laelapa/CompanyRegistry/internal/config"
	"github.com/Laelapa/CompanyRegistry/logging"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("FATAL: %v\n", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger, err := logging.NewLogger(cfg.Logging)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			log.Printf("WARNING: failed to sync logger: %v", syncErr)
		}
	}()

	dbPool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbPool.Close()
	logger.Info("Database pool initialized")

	// Verify db connection
	if dbPingErr := dbPool.Ping(ctx); dbPingErr != nil {
		return fmt.Errorf("failed to ping database: %w", dbPingErr)
	}
	logger.Info("Database connection verified")

	return nil
}
