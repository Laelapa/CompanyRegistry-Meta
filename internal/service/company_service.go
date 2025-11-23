package service

import (
	"context"
	"fmt"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"github.com/google/uuid"
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
}

func NewCompanyService(
	repo CompanyRepository,
	tokenAuthority *tokenauthority.TokenAuthority,
) *CompanyService {
	return &CompanyService{
		repo:           repo,
		tokenAuthority: tokenAuthority,
	}
}

// GetByName retrieves a company by its name.
// It returns domain.ErrNotFound if the company does not exist.
func (s *CompanyService) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	return s.repo.GetByName(ctx, name)
}

// Create creates a new company.
// If uniqueness constraints are violated, it returns domain.ErrConflict.
func (s *CompanyService) Create(ctx context.Context, c *domain.Company) (*domain.Company, error) {
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
	return s.repo.Create(ctx, c)
}

// Update updates an existing company.
// If the company does not exist, it returns domain.ErrNotFound.
// If uniqueness constraints are violated, it returns domain.ErrConflict.
func (s *CompanyService) Update(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	if c.ID == nil {
		return nil, fmt.Errorf("company ID is required: %w", domain.ErrBadRequest)
	}
	if c.UpdatedBy == nil {
		return nil, fmt.Errorf("updated_by is required: %w", domain.ErrBadRequest)
	}

	return s.repo.Update(ctx, c)
}

// Delete deletes a company by ID.
// It returns domain.ErrNotFound if the company does not exist.
func (s *CompanyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
