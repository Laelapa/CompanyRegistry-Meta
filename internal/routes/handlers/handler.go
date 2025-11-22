package handlers

import (
	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/service"
	"github.com/Laelapa/CompanyRegistry/logging"

	"github.com/go-playground/validator/v10"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Handler struct {
	logger         *logging.Logger
	service        *service.Service
	tokenAuthority *tokenauthority.TokenAuthority
	kafkaClient    *kgo.Client
	validator      *validator.Validate
}

func New(
	logger *logging.Logger,
	service *service.Service,
	tokenAuthority *tokenauthority.TokenAuthority,
	kafkaClient *kgo.Client,
) *Handler {
	return &Handler{
		logger:         logger,
		service:        service,
		tokenAuthority: tokenAuthority,
		kafkaClient:    kafkaClient,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}
}
