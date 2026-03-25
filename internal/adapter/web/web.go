package web

import (
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/router"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
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
