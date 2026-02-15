package router

import (
	"fmt"
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/middleware"
)

type RouterOptions struct {
}

func Router(opts RouterOptions) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})

	handler := http.Handler(router)
	handler = middleware.NoCache(handler)
	handler = middleware.Logger()(handler)
	handler = middleware.RealIP()(handler)
	handler = middleware.Recovery()(handler)

	return handler
}
