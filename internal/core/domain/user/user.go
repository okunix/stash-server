package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"string"`
	PasswordHash string     `json:"-"`
	Locked       bool       `json:"locked"`
	ExpiredAt    *time.Time `json:"expired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type GetUserParams struct {
	Username     string
	PasswordHash string
}

type UpdateUserParams struct {
	ID           uuid.UUID
	PasswordHash string
	Locked       bool
	ExpiredAt    *time.Time
}

type ListUsersParams struct {
	Limit  uint
	Offset uint
}

type AddUserParams struct {
	Username     string
	PasswordHash string
}

var (
	ErrNotFound = errors.New("user not found")
)
