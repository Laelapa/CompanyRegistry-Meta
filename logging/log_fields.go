package logging

import (
	"net/http"

	"go.uber.org/zap"
)

const (

	// Service related fields --------------------------

	// FieldService is the service identifier
	FieldService = "service"
	// FieldEnv is the environment the service is running in / the logger is configured for
	FieldEnv = "env"
	// FieldLoggingLevel is the logging level the logger is configured for
	FieldLoggingLevel = "logging_level"

	// Request related fields --------------------------

	// FieldRemoteAddr is the IP address of the client making the request
	FieldRemoteAddr = "remote_addr"
	// FieldMethod is the HTTP method of the request (GET, POST, etc.)
	FieldMethod = "method"
	// FieldPath is the URL path of the request
	FieldPath = "path"
	// FieldReferer is the Referer header from the request
	FieldReferer = "referer"

	// Other common fields -----------------------------

	FieldError = "error"
)

func (l *Logger) ReqFields(r *http.Request) []zap.Field {
	if r == nil {
		return []zap.Field{
			zap.String(FieldError, "request is nil"),
		}
	}

	return []zap.Field{
		zap.String(FieldRemoteAddr, l.FiletLogValue(getClientIP(r))),
		zap.String(FieldMethod, r.Method),
		zap.String(FieldPath, l.FiletLogValue(r.URL.Path)),
		zap.String(FieldReferer, l.FiletLogValue(r.Referer())),
	}
}
