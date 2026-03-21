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
	DB           *sql.DB
	UserService  ports.UserService
	StashService ports.StashService
}

func Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", newV1Router(opts)))

	handler := http.Handler(router)
	handler = middleware.NoCache(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RealIP(handler)
	handler = middleware.RequestID(handler)
	//handler = middleware.Recovery(handler)

	return handler
}

func newV1Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		var id string
		opts.DB.QueryRowContext(r.Context(), "SELECT gen_random_uuid();").Scan(&id)
		fmt.Fprintf(w, "%s\n", id)
	})

	router.Handle("POST /login", handlers.Login(opts.UserService).Unwrap())
	router.Handle("POST /signup", handlers.CreateUser(opts.UserService).Unwrap())

	router.Handle(
		"GET /stashes",
		middleware.Authenticated(handlers.ListStashes(opts.StashService).Unwrap()),
	)
	router.Handle(
		"POST /stashes",
		middleware.Authenticated(handlers.CreateStash(opts.StashService).Unwrap()),
	)
	router.Handle(
		"GET /stashes/{stash_id}",
		middleware.Authenticated(handlers.GetStashByID(opts.StashService).Unwrap()),
	)
	router.Handle(
		"DELETE /stashes/{stash_id}",
		middleware.Authenticated(handlers.DeleteStash(opts.StashService).Unwrap()),
	)

	router.Handle(
		"POST /stashes/{stash_id}/unlock",
		middleware.Authenticated(handlers.UnlockStash(opts.StashService).Unwrap()),
	)
	router.Handle(
		"POST /stashes/{stash_id}/lock",
		middleware.Authenticated(handlers.LockStash(opts.StashService).Unwrap()),
	)

	router.Handle(
		"GET /stashes/{stash_id}/secrets",
		middleware.Authenticated(handlers.GetSecrets(opts.StashService).Unwrap()),
	)
	router.Handle(
		"GET /stashes/{stash_id}/secrets/{entry_name}",
		middleware.Authenticated(handlers.GetSecretsEntry(opts.StashService).Unwrap()),
	)
	router.Handle(
		"PUT /stashes/{stash_id}/secrets",
		middleware.Authenticated(handlers.AddSecretsEntry(opts.StashService).Unwrap()),
	)
	router.Handle(
		"DELETE /stashes/{stash_id}/secrets/{entry_name}",
		middleware.Authenticated(handlers.RemoveSecretsEntry(opts.StashService).Unwrap()),
	)

	router.Handle(
		"GET /stashes/{stash_id}/members",
		middleware.Authenticated(handlers.GetStashMembers(opts.StashService).Unwrap()),
	)

	return router
}
