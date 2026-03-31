package middleware

import (
	"net/http"
	"strings"

	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/core/auth"
	"github.com/okunix/stash-server/internal/core/domain/user"
)

const (
	authHeaderKey    = "Authorization"
	authHeaderPrefix = "Bearer "
)

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.UserFromContext(r.Context())
		if !ok {
			jsonutil.SendMessage(w, jsonutil.Unauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := auth.UserFromContext(r.Context())
		if !ok {
			jsonutil.SendMessage(w, jsonutil.Unauthorized)
			return
		}
		if currentUser.Role != user.RoleAdmin {
			jsonutil.SendMessage(w, jsonutil.Forbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AssignUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authHeaderKey)

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString, found := strings.CutPrefix(authHeader, authHeaderPrefix)
		if !found {
			next.ServeHTTP(w, r)
			return
		}
		claims, err := auth.ParseJWT(tokenString)
		if err != nil || claims == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := auth.WithUser(r.Context(), &claims.CurrentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
