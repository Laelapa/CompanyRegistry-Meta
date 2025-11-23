package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/config"
	"github.com/Laelapa/CompanyRegistry/internal/middleware"
	"github.com/Laelapa/CompanyRegistry/internal/routes"
	"github.com/Laelapa/CompanyRegistry/internal/service"
	"github.com/Laelapa/CompanyRegistry/logging"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type App struct {
	server       *http.Server
	serverConfig *config.ServerConfig
	logger       *logging.Logger
}

func New(
	serverConfig *config.ServerConfig,
	logger *logging.Logger,
	service *service.Service,
	tokenAuthority *tokenauthority.TokenAuthority,
	kafkaClient *kgo.Client,
) *App {
	return &App{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", serverConfig.Port),
			Handler: newMux(
				// serverConfig.StaticDir,
				logger,
				service,
				tokenAuthority,
				kafkaClient,
			),
			ReadHeaderTimeout: serverConfig.Timeouts.ReadHeaderTimeout,
			ReadTimeout:       serverConfig.Timeouts.ReadTimeout,
			WriteTimeout:      serverConfig.Timeouts.WriteTimeout,
			IdleTimeout:       serverConfig.Timeouts.IdleTimeout,
		},
		serverConfig: serverConfig,
		logger:       logger,
	}
}

func newMux(
	// staticDir string,
	logger *logging.Logger,
	service *service.Service,
	tokenAuthority *tokenauthority.TokenAuthority,
	kafkaClient *kgo.Client,
) http.Handler {
	mux := routes.Setup(
		// staticDir,
		logger,
		service,
		tokenAuthority,
		kafkaClient,
	)
	return attachBasicMiddleware(mux, logger)
}

func attachBasicMiddleware(handler http.Handler, logger *logging.Logger) http.Handler {
	handler = middleware.RequestLogger(handler, logger)

	return handler
}

func (app *App) LaunchServer(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		app.logger.Info(
			"HTTP Server starting",
			zap.String(logging.FieldServerAddr, app.server.Addr),
		)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error(
				"ListenAndServe returned with error",
				zap.Error(err),
			)
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		app.logger.Info("Shutting down HTTP server...")
		app.ShutdownServer()
		return nil
	}
}

func (app *App) ShutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), app.serverConfig.ShutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error(
			"HTTP server shutdown returned with error",
			zap.Error(err),
		)
		app.logger.Warn(
			"Forcing HTTP server shutdown",
		)
		if closeErr := app.server.Close(); closeErr != nil {
			app.logger.Error(
				"HTTP server FORCED shutdown returned with error",
				zap.Error(closeErr),
			)
			return
		}
	}
	app.logger.Info("HTTP server shut down successfully")
}
