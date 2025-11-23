package logging

import (
	"go.uber.org/zap"

	"github.com/Laelapa/CompanyRegistry/internal/config"
)

type Logger struct {
	*zap.Logger
	maxHeaderLength int
}

func NewLogger(cfg config.LoggingConfig) (*Logger, error) {
	var zapConfig zap.Config

	switch cfg.LoggerSetup {
	// case "development", "dev":
	// 	zapConfig = setupDevConfig()
	// case "testing", "test":
	// 	zapConfig = setupTestConfig()
	// case "staging", "stg":
	// 	zapConfig = setupStagingConfig(cfg)
	case "production", "prod":
		fallthrough
	default:
		zapConfig = setupProdConfig(cfg)
	}

	logger, err := zapConfig.Build(zap.AddCallerSkip(1)) // Skip my wrapper
	if err != nil {
		return nil, err
	}

	// TODO: implement standard fields
	logger.Info(
		"Logger initialized",
		zap.String(FieldEnv, "production"),
		zap.String(FieldLoggingLevel, zapConfig.Level.String()),
	)

	return &Logger{logger, cfg.MaxHeaderLength}, nil
}

func setupProdConfig(cfg config.LoggingConfig) zap.Config {
	config := zap.NewProductionConfig()
	config.InitialFields = map[string]any{
		FieldService: cfg.ServiceName,
	}
	return config
}
