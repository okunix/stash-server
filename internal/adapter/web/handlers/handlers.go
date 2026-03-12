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
		message := jsonutil.WithDetail(jsonutil.InternalServerError, svcErr.AppError())
		switch svcErr.ServiceError() {
		case ports.ErrBadRequest:
			message = jsonutil.WithDetail(jsonutil.BadRequest, svcErr.AppError())
		case ports.ErrValidationError:
			message = jsonutil.WithDetail(jsonutil.ValidationError, svcErr.AppError())
		case ports.ErrNotFound:
			message = jsonutil.WithDetail(jsonutil.NotFound, svcErr.AppError())
		case ports.ErrUnauthorized:
			message = jsonutil.WithDetail(jsonutil.Unauthorized, svcErr.AppError())
		case ports.ErrForbidden:
			message = jsonutil.WithDetail(jsonutil.Forbidden, svcErr.AppError())
		case ports.ErrInternalError:
			message = jsonutil.WithDetail(jsonutil.InternalServerError, svcErr.AppError())
		}
		return message, true
	}
	return jsonutil.NewMessage(
		http.StatusInternalServerError,
		"Unknown Error",
		err.Error(),
	), false
}
