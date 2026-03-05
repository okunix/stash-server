package middleware

import (
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/webutil"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, webutil.WithRequestID(r))
	})
}
