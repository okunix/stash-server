package router

import (
	"fmt"
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/data"
	"gitlab.com/stash-password-manager/stash-server/middleware"
)

type RouterOptions struct {
}

func Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		db := data.SQLite()
		var id string
		db.QueryRowContext(r.Context(), "SELECT generate_uuid();").Scan(&id)
		fmt.Fprintf(w, "%s\n", id)
	})

	handler := http.Handler(router)
	handler = middleware.NoCache(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RealIP(handler)
	handler = middleware.Recovery(handler)

	return handler
}
