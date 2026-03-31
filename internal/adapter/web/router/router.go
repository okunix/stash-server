package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/okunix/stash-server/internal/adapter/web/handlers"
	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/adapter/web/middleware"
	"github.com/okunix/stash-server/internal/core/ports"
)

type RouterOptions struct {
	UserService  ports.UserService
	StashService ports.StashService
}

func Router(opts RouterOptions) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.Recovery,
		chiMiddleware.StripSlashes,
		chiMiddleware.CleanPath,
		middleware.AssignUser,
		chiMiddleware.RequestID,
		middleware.RealIP,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
			MaxAge:         300,
		}),
		middleware.Logger,
		middleware.NoCache,
	)

	router.Mount("/api/v1/", newV1Router(opts))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		jsonutil.SendMessage(w, jsonutil.NotFound)
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		jsonutil.SendMessage(w, jsonutil.MethodNotAllowed)
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
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
