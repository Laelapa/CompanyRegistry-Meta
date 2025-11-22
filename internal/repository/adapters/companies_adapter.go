package adapters

import (
	"context"
	"errors"

	"github.com/Laelapa/CompanyRegistry/internal/domain"
	"github.com/Laelapa/CompanyRegistry/internal/repository"
	"github.com/Laelapa/CompanyRegistry/util/typeconvert"

	"github.com/jackc/pgx/v5/pgtype"
)

type PGCompanyRepoAdapter struct {
	q *repository.Queries
}

func NewPGCompanyRepoAdapter(q *repository.Queries) *PGCompanyRepoAdapter {
	return &PGCompanyRepoAdapter{q: q}
}

func (p *PGCompanyRepoAdapter) Create(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	if c.Name == nil {
		return nil, errors.New("company name is required")
	}
	if c.EmployeeCount == nil {
		return nil, errors.New("employee count is required")
	}
	if c.Registered == nil {
		return nil, errors.New("registered status is required")
	}
	if c.CompanyType == nil {
		return nil, errors.New("company type is required")
	}

	params := repository.CreateCompanyParams{
		Name:          *c.Name,
		EmployeeCount: *c.EmployeeCount,
		Registered:    *c.Registered,
		CompanyType:   string(*c.CompanyType),
	}

	if c.Description != nil {
		params.Description = pgtype.Text{String: *c.Description, Valid: true}
	}
	if c.CreatedBy != nil {
		params.CreatedBy = typeconvert.GoogleUUIDToPgtypeUUID(*c.CreatedBy)
	}

	dbCompany, err := p.q.CreateCompany(ctx, params)
	if err != nil {
		return nil, err
	}

	return p.toDomainType(&dbCompany), nil
}

func (p *PGCompanyRepoAdapter) toDomainType(c *repository.Company) *domain.Company {
	ct := domain.CompanyType(c.CompanyType)
	cb := typeconvert.PgtypeUUIDToGoogleUUID(c.CreatedBy)
	ub := typeconvert.PgtypeUUIDToGoogleUUID(c.UpdatedBy)

	return &domain.Company{
		ID:            &c.ID,
		Name:          &c.Name,
		Description:   typeconvert.PgtypeTextToPtrString(c.Description),
		EmployeeCount: &c.EmployeeCount,
		Registered:    &c.Registered,
		CompanyType:   &ct,
		CreatedBy:     &cb,
		UpdatedBy:     &ub,
	}
}
