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
	Detail  any    `json:"detail,omitempty"`
}

func NewMessage(code int, message string, detail ...any) Message {
	var messageDetail any
	if len(detail) > 0 {
		messageDetail = detail[0]
	}
	return Message{Code: code, Message: message, Detail: messageDetail}
}

func SendMessage(w http.ResponseWriter, m Message) error {
	return Write(w, m.Code, m)
}

var (
	Ok                  = NewMessage(http.StatusOK, "Ok")
	Created             = NewMessage(http.StatusCreated, "Created")
	NotFound            = NewMessage(http.StatusNotFound, "Not Found")
	Forbidden           = NewMessage(http.StatusForbidden, "Forbidden")
	BadRequest          = NewMessage(http.StatusBadRequest, "Bad Request")
	Unauthorized        = NewMessage(http.StatusUnauthorized, "Unauthorized")
	TooManyRequests     = NewMessage(http.StatusTooManyRequests, "Too Many Requests")
	InternalServerError = NewMessage(http.StatusInternalServerError, "Internal Server Error")
	ValidationError     = NewMessage(http.StatusBadRequest, "Validation Error")
)

func WithDetail(m Message, detail any) Message {
	var d any
	if err, ok := detail.(error); ok {
		raw := err.Error()
		var parsed any
		if json.Unmarshal([]byte(raw), &parsed) == nil {
			d = parsed
		} else {
			d = raw
		}
	}
	m.Detail = d
	return m
}
