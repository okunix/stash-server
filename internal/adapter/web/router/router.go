package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/handlers"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/middleware"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type RouterOptions struct {
	DB          *sql.DB
	UserService ports.UserService
}

func Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		var id string
		opts.DB.QueryRowContext(r.Context(), "SELECT gen_random_uuid();").Scan(&id)
		fmt.Fprintf(w, "%s\n", id)
	})

	router.Handle(
		"GET /authcheck",
		middleware.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("you authenticated"))
		})),
	)

	router.Handle("POST /login", handlers.Login(opts.UserService).Unwrap())
	router.Handle("POST /signup", handlers.CreateUser(opts.UserService).Unwrap())

	handler := http.Handler(router)
	//handler = middleware.Authenticated(handler)
	handler = middleware.NoCache(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RealIP(handler)
	handler = middleware.RequestID(handler)
	//handler = middleware.Recovery(handler)

	return handler
}
