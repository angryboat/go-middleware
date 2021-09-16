package middleware

import (
	"net"
	"net/http"
	"strings"
)

// RemoteIP will return the IP address of the request. It searches the request
// headers X-Forwarded-For then X-Original-Forwarded-For. If neither headers
// contain IP addresses the request RemoteAddr will be used. If the RemoteAddr
// can't be parsed, a fallback will be returned "127.0.0.1".
func RemoteIP(r *http.Request) []string {
	if val := r.Header.Get(`X-Forwarded-For`); val != "" {
		return strings.Split(val, ",")
	}
	if val := r.Header.Get(`X-Original-Forwarded-For`); val != "" {
		return strings.Split(val, ",")
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	if ip == "" {
		ip = "127.0.0.1"
	}

	return []string{ip}
}
