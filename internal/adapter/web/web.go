package web

import (
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/router"
)

type ServerOptions struct {
	Addr string
}

func NewServer(opts ServerOptions) *http.Server {
	handler := router.Router(router.RouterOptions{})
	return &http.Server{
		Addr:           opts.Addr,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
	}
}
