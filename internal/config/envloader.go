package config

import (
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
	Timeouts        ServerTimeoutsConfig
}

type ServerTimeoutsConfig struct {
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
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
	// Server Timeouts
	defaultReadHeaderTimeout = 10 * time.Second  // For slow headers
	defaultReadTimeout       = 30 * time.Second  // For slow requests
	defaultWriteTimeout      = 30 * time.Second  // For preventing clients from keeping connections open
	defaultIdleTimeout       = 120 * time.Second // For closing idle connections

	// Auth
	defaultMaxHeaderLength = 1024
	defaultJwtLifetime     = 15 * time.Minute
	defaultJwtSecret       = "Default JWT Secret - DO NOT USE IN PRODUCTION" //nolint:gosec // hardcoded secret for dev/testing

	// Logging
	defaultServiceName = "my-service"
	defaultLoggerSetup = defaultEnv
)

// Load reads configuration from environment variables, doing validation and applying defaults.
// It returns an initialized Config struct or an error if required variables are missing or invalid.
// The Config struct should, ideally, be treated as readonly after loading.
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
			Port:            getEnvWithFallbackAndCustomValidation("SERVER_PORT", "8080", validatePort),
			ShutdownTimeout: getEnvDurationWithFallback("SERVER_SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
			StaticDir:       getEnvWithFallback("SERVER_STATIC_DIR", defaultStaticDir),
			Timeouts: ServerTimeoutsConfig{
				ReadHeaderTimeout: getEnvDurationWithFallback("SERVER_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout),
				ReadTimeout:       getEnvDurationWithFallback("SERVER_READ_TIMEOUT", defaultReadTimeout),
				WriteTimeout:      getEnvDurationWithFallback("SERVER_WRITE_TIMEOUT", defaultWriteTimeout),
				IdleTimeout:       getEnvDurationWithFallback("SERVER_IDLE_TIMEOUT", defaultIdleTimeout),
			},
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
