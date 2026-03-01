package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/middleware"
)

type RouterOptions struct {
	DB *sql.DB
}

func Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		var id string
		opts.DB.QueryRowContext(r.Context(), "SELECT gen_random_uuid();").Scan(&id)
		fmt.Fprintf(w, "%s\n", id)
	})

	handler := http.Handler(router)
	handler = middleware.NoCache(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RealIP(handler)
	//handler = middleware.Recovery(handler)

	return handler
}
