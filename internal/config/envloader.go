package config

import (
	"fmt"
	"log" //nolint:depguard // using standard log for config warnings before structured logging is initialized
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	DB          DatabaseConfig
	Auth        AuthConfig
	Kafka       KafkaConfig
	Logging     LoggingConfig
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	StaticDir       string
}

type DatabaseConfig struct {
	URL string
}

type AuthConfig struct {
	JwtSecret   string
	JwtIssuer   string
	JwtLifetime time.Duration
}

type KafkaConfig struct {
	ClientID string
	Brokers  []string
	Topic    struct {
		Mutations string
	}
}

type LoggingConfig struct {
	ServiceName     string
	LoggerSetup     string
	MaxHeaderLength int // where to truncate long headers in logs
}

const (
	// General
	defaultEnv = "prod"

	// Server
	defaultShutdownTimeout = 5 * time.Second
	defaultStaticDir       = "./static"

	// Auth
	defaultMaxHeaderLength = 1024
	defaultJwtLifetime     = 15 * time.Minute
	defaultJwtSecret       = "Default JWT Secret - DO NOT USE IN PRODUCTION" //nolint:gosec // hardcoded secret for dev/testing

	// Logging
	defaultServiceName = "my-service"
	defaultLoggerSetup = defaultEnv
)

func getEnvWithFallback(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: environment variable %v not set, falling back to %v", key, fallback)
		return fallback
	}
	return val
}

func getEnvRequired(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %v is required but not set", key)
	}
	return val, nil
}

func getEnvWithFallbackAndValidOptions(key string, fallback string, validoptions ...string) string {
	val := os.Getenv(key)
	if !slices.Contains(validoptions, val) {
		log.Printf("WARNING: invalid value %v for env %v, falling back to %v", val, key, fallback)
		return fallback
	}
	return val
}

func getEnvIntWithFallback(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: env %v not set, falling back to %v", key, fallback)
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("WARNING: could not parse int from env %v, falling back to %v", key, fallback)
		return fallback
	}

	return intVal
}

// getEnvDurationWithFallback retrieves a duration from an environment variable,
// uses fallback if non-positive, not parseable, or not set.
func getEnvDurationWithFallback(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: env %v not set, falling back to %v", key, fallback)
		return fallback
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("WARNING: env %v is not a valid duration format (%v), falling back to %v", key, val, fallback)
		return fallback
	}

	if dur <= 0 {
		log.Printf("WARNING: env %v must be positive, got %v, falling back to %v", key, dur, fallback)
		return fallback
	}

	return dur
}

func Load() (*Config, error) {
	validEnvs := []string{
		"prod", "production",
		"dev", "development",
		"test", "testing",
		"stg", "staging",
	}

	// Load required env variables first
	dbURL, err := getEnvRequired("DB_URL")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Environment: getEnvWithFallbackAndValidOptions("ENVIRONMENT", defaultEnv, validEnvs...),
		Server: ServerConfig{
			Port:            getEnvWithFallback("SERVER_PORT", "8080"),
			ShutdownTimeout: getEnvDurationWithFallback("SERVER_SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
			StaticDir:       getEnvWithFallback("SERVER_STATIC_DIR", defaultStaticDir),
		},
		DB: DatabaseConfig{
			URL: dbURL,
		},
		Auth: AuthConfig{
			JwtSecret:   getEnvWithFallback("JWT_SECRET", defaultJwtSecret),
			JwtIssuer:   getEnvWithFallback("JWT_ISSUER", getEnvWithFallback("SERVICE_NAME", "my-service")),
			JwtLifetime: getEnvDurationWithFallback("JWT_LIFETIME", defaultJwtLifetime),
		},
		Kafka: KafkaConfig{
			ClientID: getEnvWithFallback("SERVICE_NAME", defaultServiceName),
			Brokers:  strings.Split(getEnvWithFallback("KAFKA_BROKERS", "localhost:9092"), ","),
			Topic: struct{ Mutations string }{
				Mutations: getEnvWithFallback("KAFKA_TOPIC_MUTATIONS", "mutations"),
			},
		},
		Logging: LoggingConfig{
			ServiceName:     getEnvWithFallback("SERVICE_NAME", defaultServiceName),
			LoggerSetup:     getEnvWithFallbackAndValidOptions("LOGGER_SETUP", defaultLoggerSetup, validEnvs...),
			MaxHeaderLength: getEnvIntWithFallback("MAX_HEADER_LENGTH", defaultMaxHeaderLength),
		},
	}
	return cfg, nil
}
