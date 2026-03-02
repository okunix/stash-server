package jsonutil

import (
	"encoding/json"
	"io"
	"net/http"
)

func Write(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func Read[T any](r io.Reader) (T, error) {
	var dest T
	err := json.NewDecoder(r).Decode(&dest)
	return dest, err
}

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewMessage(code int, message string) Message {
	return Message{Code: code, Message: message}
}

func SendMessage(w http.ResponseWriter, m Message) error {
	return Write(w, m.Code, m.Message)
}

var (
	OkMessage                  = NewMessage(http.StatusOK, "ok")
	CreatedMessage             = NewMessage(http.StatusCreated, "created")
	ForbiddenMessage           = NewMessage(http.StatusForbidden, "forbidden")
	UnauthorizedMessage        = NewMessage(http.StatusUnauthorized, "unauthorized")
	TooManyRequestsMessage     = NewMessage(http.StatusTooManyRequests, "too many requests")
	InternalServerErrorMessage = NewMessage(http.StatusInternalServerError, "internal server error")
)
