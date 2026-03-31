package web

import (
	"net/http"

	"github.com/okunix/stash-server/internal/adapter/web/router"
	"github.com/okunix/stash-server/internal/core/ports"
)

type ServerOptions struct {
	Addr         string
	UserService  ports.UserService
	StashService ports.StashService
}

func NewServer(opts ServerOptions) *http.Server {
	handler := router.Router(
		router.RouterOptions{
			UserService:  opts.UserService,
			StashService: opts.StashService,
		},
	)
	return &http.Server{
		Addr:           opts.Addr,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
	}
}
