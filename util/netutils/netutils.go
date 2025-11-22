package netutils

import (
	"errors"
	"net"
	"net/http"
)

// GetClientIP extracts the client IP address from the HTTP request.
// It first checks the X-Forwarded-For header, falling back to RemoteAddr if not present.
func GetClientIP(r *http.Request) string {
	if clientIP := r.Header.Get("X-Forwarded-For"); clientIP != "" {
		if ip := net.ParseIP(clientIP); ip != nil {
			return ip.String()
		}
		return ""
	}
	return "X-Forwarded-For empty, falling back to r.RemoteAddr:" + StripPort(r.RemoteAddr)
}

// StripPort removes the port from an IP address string, if present.
// If no port is found, it returns the original string.
func StripPort(ipAddr string) string {
	host, _, err := net.SplitHostPort(ipAddr)
	if err != nil {
		return ipAddr
	}
	return host
}

func StripBearer(authHeader string) (string, error) {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return "", errors.New("invalid authorization header format")
}
