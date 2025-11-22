package middleware

import (
	"net/http"

	"github.com/Laelapa/CompanyRegistry/logging"
)

func RequestLogger(next http.Handler, logger *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Incoming request", logger.ReqFields(r)...)
		next.ServeHTTP(w, r)
	})
}
