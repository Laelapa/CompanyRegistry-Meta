//nolint:dupl // TODO: Refactor into generic HandleAuth called by both Signup and Login
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"go.uber.org/zap"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var rBody UserLoginRequest
	h.logger.Info("Processing login request", h.logger.ReqFields(r)...)

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.Warn(
			"Failed to decode login request body",
			append(h.logger.ReqFields(r), zap.Error(err))...,
		)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate contents
	if err := h.validator.Struct(rBody); err != nil {
		h.logger.Warn(
			"Invalid login request data",
			append(h.logger.ReqFields(r), zap.Error(err))...,
		)
		http.Error(w, "Bad request: Invalid data", http.StatusBadRequest)
		return
	}

	accessToken, err := h.service.User.Login(r.Context(), rBody.Username, rBody.Password)
	if err != nil {
		h.logger.Error(
			"User login failed",
			append(h.logger.ReqFields(r), zap.Error(err))...,
		)
		if errors.Is(err, domain.ErrBadCredentials) {
			http.Error(w, "Unauthorized: Login Failed", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := AuthResponse{AccessToken: accessToken}

	respMarshalled, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("Failed to marshal login response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respMarshalled); err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
	}
	h.logger.Info("Login request processed", h.logger.ReqFields(r)...)
}
