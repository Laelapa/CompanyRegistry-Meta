package handlers

import (
	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"github.com/google/uuid"
)

type CreateCompanyRequest struct {
	Name          string  `json:"name"           validate:"required,max=15"`
	Description   *string `json:"description"    validate:"omitempty,max=3000"`
	EmployeeCount *int32  `json:"employee_count" validate:"required,gte=0"`
	Registered    *bool   `json:"registered"     validate:"required"`
	CompanyType   string  `json:"company_type"   validate:"required,oneof='Corporation' 'NonProfit' 'Cooperative' 'Sole Proprietorship'"`
}

type UpdateCompanyRequest struct {
	Name          *string `json:"name"           validate:"omitempty,max=15"`
	Description   *string `json:"description"    validate:"omitempty,max=3000"`
	EmployeeCount *int32  `json:"employee_count" validate:"omitempty,gte=0"`
	Registered    *bool   `json:"registered"     validate:"omitempty"`
	CompanyType   *string `json:"company_type"   validate:"omitempty,oneof='Corporation' 'NonProfit' 'Cooperative' 'Sole Proprietorship'"`
}

type CompanyResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	EmployeeCount int32     `json:"employee_count"`
	Registered    bool      `json:"registered"`
	CompanyType   string    `json:"company_type"`
}

func convertToCompanyResponse(c *domain.Company) CompanyResponse {
	id := uuid.Nil
	if c.ID != nil {
		id = *c.ID
	}
	name := ""
	if c.Name != nil {
		name = *c.Name
	}
	empCount := int32(0)
	if c.EmployeeCount != nil {
		empCount = *c.EmployeeCount
	}
	registered := false
	if c.Registered != nil {
		registered = *c.Registered
	}
	cType := ""
	if c.CompanyType != nil {
		cType = string(*c.CompanyType)
	}
	return CompanyResponse{
		ID:            id,
		Name:          name,
		Description:   c.Description,
		EmployeeCount: empCount,
		Registered:    registered,
		CompanyType:   cType,
	}
}
