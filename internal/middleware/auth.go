package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/logging"
	"github.com/Laelapa/CompanyRegistry/util/ctxutils"
	"github.com/Laelapa/CompanyRegistry/util/netutils"
)

func AuthenticateWithJWT(
	tokenAuthority *tokenauthority.TokenAuthority,
	logger *logging.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Unauthorized request: Missing Authorization header", logger.ReqFields(r)...)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Extract JWT from bearer scheme
			tokenString, err := netutils.StripBearer(authHeader)
			if err != nil {
				logger.Warn(
					"Unauthorized request: Invalid Authorization header",
					append(logger.ReqFields(r), zap.Error(err))...,
				)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate JWT & extract userID
			userUUID, err := tokenAuthority.ValidateJWT(tokenString)
			if err != nil {
				logger.Warn(
					"Unauthorized request: Invalid token",
					append(logger.ReqFields(r), zap.Error(err))...,
				)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := ctxutils.SetUserIDInContext(r.Context(), userUUID)

			logger.Info(
				"Request authenticated with JWT",
				append(
					logger.ReqFields(r),
					zap.String(logging.FieldUserID, userUUID.String()),
				)...,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
