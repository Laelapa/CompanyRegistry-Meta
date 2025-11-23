package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"go.uber.org/zap"
)

// HandleGetCompanyByName processes requests to retrieve a company by name.
func (h *Handler) HandleGetCompanyByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.logger.Warn("Company name missing from path", h.logger.ReqFields(r)...)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	h.logger.Info("Processing Get Company By Name request", h.logger.ReqFields(r)...)

	company, err := h.service.Company.GetByName(r.Context(), name)
	if err != nil {
		h.logger.Error("Failed to get company", append(h.logger.ReqFields(r), zap.Error(err))...)
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := convertToCompanyResponse(company)

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
	h.logger.Info("Company Get request processed", h.logger.ReqFields(r)...)
}
