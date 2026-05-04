package httputil

import (
	"net"
	"net/http"
)

// GetClientIP extracts the client IP address from the request.
// It first checks the X-Real-IP header and falls back to RemoteAddr.
func GetClientIP(r *http.Request) string {
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
