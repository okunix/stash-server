package router

import (
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/handlers"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/middleware"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type RouterOptions struct {
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
	handler = middleware.AddTrailingSlash(handler)
	handler = middleware.Recovery(handler)

	return handler
}

func newV1StashRouter(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /{$}",
		handlers.ListStashes(opts.StashService).Unwrap())

	router.Handle("POST /{$}",
		handlers.CreateStash(opts.StashService).Unwrap())

	router.Handle("GET /{stash_id}",
		handlers.GetStashByID(opts.StashService).Unwrap())

	router.Handle("DELETE /{stash_id}",
		handlers.DeleteStash(opts.StashService).Unwrap())

	router.Handle("POST /{stash_id}/unlock",
		handlers.UnlockStash(opts.StashService).Unwrap())

	router.Handle("POST /{stash_id}/lock",
		handlers.LockStash(opts.StashService).Unwrap())

	router.Handle("GET /{stash_id}/secrets",
		handlers.GetSecrets(opts.StashService).Unwrap())

	router.Handle("GET /{stash_id}/secrets/{entry_name}",
		handlers.GetSecretsEntry(opts.StashService).Unwrap())

	router.Handle("PUT /{stash_id}/secrets",
		handlers.AddSecretsEntry(opts.StashService).Unwrap())

	router.Handle("DELETE /{stash_id}/secrets/{entry_name}",
		handlers.RemoveSecretsEntry(opts.StashService).Unwrap())

	router.Handle("GET /{stash_id}/members",
		handlers.GetStashMembers(opts.StashService).Unwrap())

	router.Handle("POST /{stash_id}/members",
		handlers.AddStashMember(opts.StashService).Unwrap())

	router.Handle("DELETE /{stash_id}/members/{user_id}",
		handlers.RemoveStashMember(opts.StashService).Unwrap())

	handler := http.Handler(router)
	handler = middleware.Authenticated(handler)

	return handler
}

func newV1UserRouter(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /{user_id}",
		handlers.GetUserByID(opts.UserService).Unwrap())

	router.Handle("POST /{$}",
		middleware.Admin(handlers.CreateUser(opts.UserService).Unwrap()))

	handler := http.Handler(router)
	handler = middleware.Authenticated(handler)

	return handler
}

func newV1Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.Handle("/stashes/", http.StripPrefix("/stashes", newV1StashRouter(opts)))
	router.Handle("/users/", http.StripPrefix("/users", newV1UserRouter(opts)))

	router.Handle("GET /whoami/",
		middleware.Authenticated(handlers.Whoami(opts.UserService).Unwrap()))

	router.Handle("POST /login/",
		handlers.Login(opts.UserService).Unwrap())

	return router
}
