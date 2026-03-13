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
		if errors.Is(err, io.EOF) {
			jsonutil.SendMessage(w, jsonutil.WithDetail(jsonutil.BadRequest, "Empty request body"))
			return
		}
		if msg, ok := fromServiceError(err); ok {
			jsonutil.SendMessage(w, msg)
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

func fromServiceError(err error) (jsonutil.Message, bool) {
	var svcErr ports.Error
	if errors.As(err, &svcErr) {
		appErr := svcErr.AppError()
		message := jsonutil.WithDetail(jsonutil.InternalServerError, appErr)
		switch svcErr.ServiceError() {
		case ports.ErrBadRequest:
			message = jsonutil.WithDetail(jsonutil.BadRequest, appErr)
		case ports.ErrValidationError:
			message = jsonutil.WithDetail(jsonutil.ValidationError, appErr)
		case ports.ErrNotFound:
			message = jsonutil.WithDetail(jsonutil.NotFound, appErr)
		case ports.ErrUnauthorized:
			message = jsonutil.WithDetail(jsonutil.Unauthorized, appErr)
		case ports.ErrForbidden:
			message = jsonutil.WithDetail(jsonutil.Forbidden, appErr)
		case ports.ErrInternalError:
			message = jsonutil.WithDetail(jsonutil.InternalServerError, appErr)
		}
		return message, true
	}
	return jsonutil.NewMessage(
		http.StatusInternalServerError,
		"Unknown Error",
		err.Error(),
	), false
}
