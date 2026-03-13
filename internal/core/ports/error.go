package ports

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ValidationError map[string]string

func (e ValidationError) Error() string {
	jsonBytes, _ := json.Marshal(e)
	return string(jsonBytes)
}

type Error struct {
	appError error
	svcError error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s",
		e.ServiceError().Error(),
		e.AppError().Error(),
	)
}

func NewError(svcError, appError error) error {
	return Error{appError, svcError}
}

func (e Error) AppError() error {
	return e.appError
}
func (e Error) ServiceError() error {
	return e.svcError
}

// service errors (svcError)
var (
	ErrBadRequest      = errors.New("bad request")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInternalError   = errors.New("internal error")
	ErrValidationError = errors.New("validation error")
)

func InternalError(err error) error {
	return NewError(ErrInternalError, err)
}

func NotFoundError(err error) error {
	return NewError(ErrNotFound, err)
}

func BadRequestError(err error) error {
	return NewError(ErrBadRequest, err)
}

func UnauthorizedError(err error) error {
	return NewError(ErrUnauthorized, err)
}

func ForbiddenError(err error) error {
	return NewError(ErrForbidden, err)
}

func NewValidationError(problems map[string]string) error {
	return NewError(ErrValidationError, ValidationError(problems))
}
