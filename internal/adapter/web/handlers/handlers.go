package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/webutil"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (f apiFunc) Unwrap() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		if errors.As(err, &ports.ValidationError{}) {
			jsonutil.Write(w, http.StatusBadRequest, err.(*ports.ValidationError))
			return
		}
		if errors.Is(err, io.EOF) {
			jsonutil.SendMessage(w, jsonutil.BadRequest)
			return
		}
		requestID, _ := webutil.RequestID(r)
		slog.Error("error occured while processing request",
			"requestID", requestID,
			"error", err.Error(),
		)
		jsonutil.SendMessage(w, jsonutil.InternalServerError)
	})
}
