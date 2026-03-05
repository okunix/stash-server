package webutil

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const requestIDKey = "requestID"

func WithRequestID(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), requestIDKey, uuid.NewString())
	return r.WithContext(ctx)
}

func RequestID(r *http.Request) (string, bool) {
	contextValue := r.Context().Value(requestIDKey)
	if contextValue == nil {
		return "", false
	}
	return contextValue.(string), true
}
