package handlers

import (
	"errors"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) HandleDeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, pErr := uuid.Parse(r.PathValue("id"))
	if pErr != nil {
		h.logger.Warn("Invalid company ID in path", append(h.logger.ReqFields(r), zap.Error(pErr))...)
		http.Error(w, "Bad request: Invalid ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Processing Delete Company request", h.logger.ReqFields(r)...)

	if err := h.service.Company.Delete(r.Context(), id); err != nil {
		h.logger.Error("Failed to delete company", append(h.logger.ReqFields(r), zap.Error(err))...)
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Info("Successfully deleted company", h.logger.ReqFields(r)...)
}
