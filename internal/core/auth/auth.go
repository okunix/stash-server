package auth

import (
	"context"

	"github.com/google/uuid"
)

const (
	userContextKey = "user"
)

type CurrentUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}

func WithUser(ctx context.Context, user *CurrentUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*CurrentUser, bool) {
	user, ok := ctx.Value(userContextKey).(*CurrentUser)
	return user, ok
}
