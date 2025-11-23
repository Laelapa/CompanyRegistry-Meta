package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/internal/domain"
	"github.com/Laelapa/CompanyRegistry/util/ctxutils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// HandleUpdateCompany processes requests to (partially) update a company.
// It expects a userID in the request context - set by the jwt authentication middleware.
func (h *Handler) HandleUpdateCompany(w http.ResponseWriter, r *http.Request) {
	id, pErr := uuid.Parse(r.PathValue("id"))
	if pErr != nil {
		h.logger.Warn("Invalid company ID in path", append(h.logger.ReqFields(r), zap.Error(pErr))...)
		http.Error(w, "Bad request: Invalid ID", http.StatusBadRequest)
		return
	}

	var rBody UpdateCompanyRequest
	h.logger.Info("Processing Update Company request", h.logger.ReqFields(r)...)

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

	var companyType *domain.CompanyType
	if rBody.CompanyType != nil {
		ct := domain.CompanyType(*rBody.CompanyType)
		companyType = &ct
	}

	uc := &domain.Company{
		ID:            &id,
		Name:          rBody.Name,
		Description:   rBody.Description,
		EmployeeCount: rBody.EmployeeCount,
		Registered:    rBody.Registered,
		CompanyType:   companyType,
		UpdatedBy:     &userID,
	}

	updatedCompany, err := h.service.Company.Update(r.Context(), uc)
	if err != nil {
		h.logger.Error("Failed to update company", append(h.logger.ReqFields(r), zap.Error(err))...)
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "Company not found", http.StatusNotFound)
		case errors.Is(err, domain.ErrConflict):
			http.Error(w, "Conflict: Company name already exists", http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := convertToCompanyResponse(updatedCompany)

	respMarshalled, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respMarshalled); err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
	}
	h.logger.Info("Company updated successfully", h.logger.ReqFields(r)...)
}
