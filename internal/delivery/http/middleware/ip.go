package middleware

import (
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "172.18.0.1" || ip == "::1" {
		return "localhost"
	}
	return ip
}
