package middleware

import (
	"net/http"
	"strings"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/core/auth"
)

const (
	authHeaderKey    = "Authorization"
	authHeaderPrefix = "Bearer "
)

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authHeaderKey)
		tokenString, found := strings.CutPrefix(authHeader, authHeaderPrefix)
		if !found {
			jsonutil.SendMessage(w, jsonutil.Unauthorized)
			return
		}
		claims, err := auth.ParseJWT(tokenString)
		if err != nil || claims == nil {
			jsonutil.SendMessage(w, jsonutil.Unauthorized)
			return
		}
		ctx := auth.WithUser(r.Context(), &claims.CurrentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
