package user

import (
	"context"
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

type Repository interface {
	ListUsers(ctx context.Context, params ListUsersParams) ([]*User, int64, error)
	GetUser(ctx context.Context, params GetUserParams) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, params UpdateUserParams) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

var (
	ErrNotFound = errors.New("user not found")
)
