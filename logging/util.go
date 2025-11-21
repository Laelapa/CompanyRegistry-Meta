package logging

import "net/http"

// getClientIP extracts the client IP address from the HTTP request.
// It first checks the X-Forwarded-For header, falling back to RemoteAddr if not present.
func getClientIP(r *http.Request) string {
	if clientIP := r.Header.Get("X-Forwarded-For"); clientIP != "" {
		return clientIP
	}
	return "X-Forwarded-For empty, falling back to r.RemoteAddr:" + r.RemoteAddr
}
