package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/user"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) AddUser(
	ctx context.Context,
	params user.AddUserParams,
) (*user.User, error) {
	panic("unimplemented")
}

func (u *userRepository) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	panic("unimplemented")
}

func (u *userRepository) GetUser(
	ctx context.Context,
	params user.GetUserParams,
) (*user.User, error) {
	panic("unimplemented")
}

func (u *userRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*user.User, error) {
	panic("unimplemented")
}

func (u *userRepository) ListUsers(
	ctx context.Context,
	params user.ListUsersParams,
) ([]*user.User, int64, error) {
	panic("unimplemented")
}

func (u *userRepository) UpdateUser(
	ctx context.Context,
	params user.UpdateUserParams,
) (*user.User, error) {
	panic("unimplemented")
}
