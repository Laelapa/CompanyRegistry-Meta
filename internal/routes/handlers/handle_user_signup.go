package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"go.uber.org/zap"
)

func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var rBody UserSignupRequest
	h.logger.Info("Processing signup request", h.logger.ReqFields(r)...)

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
		h.logger.Warn(
			"Failed to decode signup request body",
			append(h.logger.ReqFields(r), zap.Error(err))...,
		)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate contents
	if err := h.validator.Struct(rBody); err != nil {
		h.logger.Warn(
			"Invalid signup request data",
			append(h.logger.ReqFields(r), zap.Error(err))...,
		)
		http.Error(w, "Bad request: Invalid data", http.StatusBadRequest)
		return
	}

	accessToken, err := h.service.User.Register(r.Context(), rBody.Username, rBody.Password)
	if err != nil {
		h.logger.Error(
			"User registration failed",
			append(h.logger.ReqFields(r), zap.Error(err))...,
		)
		switch {
		case errors.Is(err, domain.ErrConflict):
			http.Error(w, "Conflict: User already exists", http.StatusConflict)
		case errors.Is(err, domain.ErrBadCredentials):
			http.Error(w, "Bad request: Invalid credentials", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := AuthResponse{AccessToken: accessToken}

	respMarshalled, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("Failed to marshal signup response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(respMarshalled); err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))
	}
	h.logger.Info("Signup request processed succesfully", h.logger.ReqFields(r)...)
}
