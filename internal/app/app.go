package app

import (
	"net/http"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/config"
	"github.com/Laelapa/CompanyRegistry/internal/repository"
	"github.com/Laelapa/CompanyRegistry/logging"
	"github.com/twmb/franz-go/pkg/kgo"
)

type App struct {
	server       *http.Server
	serverConfig *config.ServerConfig
	logger       *logging.Logger
}

func New(
	serverConfig *config.ServerConfig,
	logger *logging.Logger,
	queries *repository.Queries,
	tokenAuthority *tokenauthority.TokenAuthority,
	kafkaClient *kgo.Client,
) *App {
	return &App{
		server: &http.Server{
			Addr:
		}
	}
}
