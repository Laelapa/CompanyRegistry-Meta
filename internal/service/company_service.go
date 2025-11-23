package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/domain"
	"github.com/Laelapa/CompanyRegistry/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CompanyRepository interface {
	Create(ctx context.Context, c *domain.Company) (*domain.Company, error)
	GetByName(ctx context.Context, name string) (*domain.Company, error)
	Update(ctx context.Context, c *domain.Company) (*domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CompanyService struct {
	repo           CompanyRepository
	tokenAuthority *tokenauthority.TokenAuthority
	logger         *logging.Logger
	producer       EventProducer
	topic          string
}

func NewCompanyService(
	repo CompanyRepository,
	tokenAuthority *tokenauthority.TokenAuthority,
	logger *logging.Logger,
	producer EventProducer,
	topic string,
) *CompanyService {
	return &CompanyService{
		repo:           repo,
		tokenAuthority: tokenAuthority,
		logger:         logger,
		producer:       producer,
		topic:          topic,
	}
}

// GetByName retrieves a company by its name.
// It returns domain.ErrNotFound if the company does not exist.
func (u *CompanyService) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	return u.repo.GetByName(ctx, name)
}

// Create creates a new company.
// If uniqueness constraints are violated, it returns domain.ErrConflict.
func (u *CompanyService) Create(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	if c.Name == nil {
		return nil, fmt.Errorf("company name is required: %w", domain.ErrBadRequest)
	}
	if c.EmployeeCount == nil {
		return nil, fmt.Errorf("employee count is required: %w", domain.ErrBadRequest)
	}
	if c.Registered == nil {
		return nil, fmt.Errorf("registered status is required: %w", domain.ErrBadRequest)
	}
	if c.CompanyType == nil {
		return nil, fmt.Errorf("company type is required: %w", domain.ErrBadRequest)
	}

	createdCompany, err := u.repo.Create(ctx, c)
	if err != nil {
		return nil, err
	}

	go u.publishEvent(context.Background(), "CREATE", *createdCompany.ID)
	return createdCompany, nil
}

// Update updates an existing company.
// If the company does not exist, it returns domain.ErrNotFound.
// If uniqueness constraints are violated, it returns domain.ErrConflict.
func (u *CompanyService) Update(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	if c.ID == nil {
		return nil, fmt.Errorf("company ID is required: %w", domain.ErrBadRequest)
	}
	if c.UpdatedBy == nil {
		return nil, fmt.Errorf("updated_by is required: %w", domain.ErrBadRequest)
	}

	updatedCompany, err := u.repo.Update(ctx, c)
	if err != nil {
		return nil, err
	}

	go u.publishEvent(context.Background(), "UPDATE", *updatedCompany.ID)
	return updatedCompany, nil
}

// Delete deletes a company by ID.
// It returns domain.ErrNotFound if the company does not exist.
func (u *CompanyService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.publishEvent(context.Background(), "DELETE", id)
	return nil
}

func (u *CompanyService) publishEvent(ctx context.Context, eventType string, id uuid.UUID) {
	// pub/sub not configured
	if u.producer == nil {
		return
	}

	eventData := map[string]any{
		"event": eventType,
	}

	marshalledEvent, err := json.Marshal(eventData)
	if err != nil {
		u.logger.Error("Failed to marshal event", zap.Error(err))
		return
	}

	if pErr := u.producer.Produce(ctx, u.topic, id.String(), marshalledEvent); pErr != nil {
		u.logger.Error("Failed to produce event", zap.Error(pErr))
		return
	}
}
