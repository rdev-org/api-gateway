package shared

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		clientIP := strings.TrimSpace(ips[0])
		if clientIP != "" {
			return clientIP
		}
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		clientIP := strings.TrimSpace(xri)
		if clientIP != "" {
			return clientIP
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
