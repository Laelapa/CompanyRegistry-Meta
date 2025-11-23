package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/internal/domain"
	"github.com/Laelapa/CompanyRegistry/util/ctxutils"

	"go.uber.org/zap"
)

// HandleCreateCompany processes requests to create a new company.
// It expects a userID in the request context - set by the jwt authentication middleware.
func (h *Handler) HandleCreateCompany(w http.ResponseWriter, r *http.Request) {
	var rBody CreateCompanyRequest
	h.logger.Info("Processing Create Company request", h.logger.ReqFields(r)...)

	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.Warn("Failed to decode request body", append(h.logger.ReqFields(r), zap.Error(err))...)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(rBody); err != nil {
		h.logger.Warn("Invalid request data", append(h.logger.ReqFields(r), zap.Error(err))...)
		http.Error(w, "Bad request: Invalid data", http.StatusBadRequest)
		return
	}

	userID, ok := ctxutils.GetUserIDFromContext(r.Context())
	if !ok {
		h.logger.Error("Failed to get user ID from context", h.logger.ReqFields(r)...)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	companyType := domain.CompanyType(rBody.CompanyType)
	newCompany := &domain.Company{
		Name:          &rBody.Name,
		Description:   rBody.Description,
		EmployeeCount: rBody.EmployeeCount,
		Registered:    rBody.Registered,
		CompanyType:   &companyType,
		CreatedBy:     &userID,
	}

	createdCompany, err := h.service.Company.Create(r.Context(), newCompany)
	if err != nil {
		h.logger.Error("Failed to create company", append(h.logger.ReqFields(r), zap.Error(err))...)
		if errors.Is(err, domain.ErrConflict) {
			http.Error(w, "Conflict: Company already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := convertToCompanyResponse(createdCompany)

	respMarshalled, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(respMarshalled); err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
	}
	h.logger.Info("Company creation request processed successfully", h.logger.ReqFields(r)...)
}
