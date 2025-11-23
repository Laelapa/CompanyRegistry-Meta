package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/app"
	"github.com/Laelapa/CompanyRegistry/internal/config"
	"github.com/Laelapa/CompanyRegistry/internal/events"
	"github.com/Laelapa/CompanyRegistry/internal/repository"
	"github.com/Laelapa/CompanyRegistry/internal/repository/adapters"
	"github.com/Laelapa/CompanyRegistry/internal/service"
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
	tokenAuthority := tokenauthority.New(&cfg.Auth)
	dbPool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbPool.Close()
	logger.Info("Database pool initialized")

	// Verify db connection
	logger.Info("Pinging database...")
	if dbPingErr := dbPool.Ping(ctx); dbPingErr != nil {
		return fmt.Errorf("failed to ping database: %w", dbPingErr)
	}
	logger.Info("Database connection verified")

	var kgoClient *kgo.Client // nil if Kafka not configured
	var producer service.EventProducer
	kafkaBrokers := cfg.Kafka.Brokers
	if len(kafkaBrokers) > 0 {
		client, kErr := kgo.NewClient(
			kgo.SeedBrokers(kafkaBrokers...),
			kgo.ClientID(cfg.Kafka.ClientID),
		)
		if kErr != nil {
			logger.Warn("Failed to initialize Kafka client, continuing without Kafka", zap.Error(kErr))
		} else {
			kgoClient = client
			producer = events.NewProducer(kgoClient)
			logger.Info(
				"Kafka client initialized",
				zap.Strings(logging.FieldKafkaBrokers, kafkaBrokers),
			)
		}
	} else {
		logger.Warn("No Kafka brokers configured, skipping Kafka client initialization")
	} // If no brokers, kafkaClient & producer remain nil

	queries := repository.New(dbPool)
	service := &service.Service{
		User: service.NewUserService(
			adapters.NewPGUserRepoAdapter(queries),
			tokenAuthority,
			logger,
			producer,
			cfg.Kafka.Topic.UserMutations,
		),
		Company: service.NewCompanyService(
			adapters.NewPGCompanyRepoAdapter(queries),
			tokenAuthority,
			logger,
			producer,
			cfg.Kafka.Topic.CompanyMutations,
		),
	}

	app := app.New(
		&cfg.Server,
		logger,
		service,
		tokenAuthority,
		kgoClient,
	)
	if err = app.LaunchServer(ctx); err != nil {
		return fmt.Errorf("failed to launch server: %w", err)
	}
	return nil
}
