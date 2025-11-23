package handlers

import "github.com/google/uuid"

type CreateCompanyRequest struct {
	Name          string  `json:"name"           validate:"required, max=15"`
	Description   *string `json:"description"    validate:"omitempty"`
	EmployeeCount *int32  `json:"employee_count" validate:"required,gte=0"`
	Registered    *bool   `json:"registered"     validate:"required"`
	CompanyType   string  `json:"company_type"   validate:"required,oneof='Corporation','NonProfit','Cooperative','Sole Proprietorship'"`
}

type UpdateCompanyRequest struct {
	Name          *string `json:"name"           validate:"omitempty, max=15"`
	Description   *string `json:"description"    validate:"omitempty"`
	EmployeeCount *int32  `json:"employee_count" validate:"omitempty"`
	Registered    *bool   `json:"registered"     validate:"omitempty"`
	CompanyType   *string `json:"company_type"   validate:"omitempty,oneof='Corporation','NonProfit','Cooperative','Sole Proprietorship'"`
}

type CompanyResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	EmployeeCount int32     `json:"employee_count"`
	Registered    bool      `json:"registered"`
	CompanyType   string    `json:"company_type"`
}
