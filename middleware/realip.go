package middleware

import (
	"net"
	"net/http"
	"strings"
)

var (
	trueClientIP  = http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

func RealIP() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ip string
			if tcip := r.Header.Get(trueClientIP); tcip != "" {
				ip = tcip
			} else if xrip := r.Header.Get(xRealIP); xrip != "" {
				ip = xrip
			} else if xff := r.Header.Get(xForwardedFor); xff != "" {
				ip, _, _ = strings.Cut(xff, ",")
			}
			if ip != "" && net.ParseIP(ip) != nil {
				r.RemoteAddr = ip
			}
			next.ServeHTTP(w, r)
		})
	}
}
