package service

import (
	"context"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	Create(ctx context.Context, c *domain.Company) (*domain.Company, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
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

// Create creates a new company.
// If uniqueness constraints are violated, it returns domain.ErrConflict.
func (s *CompanyService) Create(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	return s.repo.Create(ctx, c)
}
