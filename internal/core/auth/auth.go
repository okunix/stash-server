package auth

import (
	"context"

	"github.com/google/uuid"
)

const (
	userKey = "user"
)

type User struct {
	UserID uuid.UUID
}

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}
