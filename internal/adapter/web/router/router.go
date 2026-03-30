package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/handlers"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/middleware"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type RouterOptions struct {
	UserService  ports.UserService
	StashService ports.StashService
}

func Router(opts RouterOptions) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.NoCache)
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.AssignUser)
	router.Use(chiMiddleware.CleanPath)
	router.Use(chiMiddleware.StripSlashes)
	router.Use(middleware.Recovery)

	router.Mount("/api/v1/", newV1Router(opts))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		jsonutil.SendMessage(w, jsonutil.NotFound)
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		jsonutil.SendMessage(w, jsonutil.MethodNotAllowed)
	})

	return router
}

func newV1StashRouter(opts RouterOptions) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Authenticated)

	router.Handle("GET /",
		handlers.ListStashes(opts.StashService).Unwrap())

	router.Handle("POST /",
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

	router.Handle("PATCH /{stash_id}",
		handlers.UpdateStash(opts.StashService).Unwrap())

	return router
}

func newV1UserRouter(opts RouterOptions) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Authenticated)

	router.Handle("GET /",
		handlers.ListUsers(opts.UserService).Unwrap())

	router.Handle("POST /",
		middleware.Admin(handlers.CreateUser(opts.UserService).Unwrap()))

	router.Handle("GET /{user_id}",
		handlers.GetUserByUsernameOrID(opts.UserService).Unwrap())

	return router
}

func newV1AuthRouter(opts RouterOptions) http.Handler {
	router := chi.NewRouter()

	router.Handle("POST /login",
		handlers.Login(opts.UserService).Unwrap())

	router.Handle("GET /whoami",
		middleware.Authenticated(handlers.Whoami(opts.UserService).Unwrap()))

	router.Handle("PATCH /change-password",
		middleware.Authenticated(handlers.ChangePassword(opts.UserService).Unwrap()))

	return router
}

func newV1Router(opts RouterOptions) http.Handler {
	router := chi.NewRouter()

	router.Mount("/users", newV1UserRouter(opts))
	router.Mount("/stashes", newV1StashRouter(opts))
	router.Mount("/auth", newV1AuthRouter(opts))

	return router
}
