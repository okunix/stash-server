package webutil

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestID(r *http.Request) (string, bool) {
	id := middleware.GetReqID(r.Context())
	return id, id != ""
}
