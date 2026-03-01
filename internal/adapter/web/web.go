package web

import (
	"database/sql"
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/router"
)

type ServerOptions struct {
	Addr string
	DB   *sql.DB
}

func NewServer(opts ServerOptions) *http.Server {
	handler := router.Router(router.RouterOptions{DB: opts.DB})
	return &http.Server{
		Addr:           opts.Addr,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
	}
}
