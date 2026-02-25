package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/domain/user"
	"gitlab.com/stash-password-manager/stash-server/sqlc"
)

type userRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

func userModelToUser(userModel *sqlc.User) *user.User {
	if userModel == nil {
		return nil
	}
	locked := userModel.Locked == 1
	return &user.User{
		ID:           userModel.ID,
		Username:     userModel.Username,
		PasswordHash: userModel.PasswordHash,
		ExpiredAt:    &userModel.ExpiredAt,
		Locked:       locked,
		CreatedAt:    userModel.CreatedAt,
	}
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db, queries: sqlc.New(db)}
}

func (u *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	rows, err := u.queries.DeleteUser(ctx, id)
	if rows == 0 {
		return user.ErrNotFound
	}
	return err
}

func (u *userRepository) GetUser(
	ctx context.Context,
	params user.GetUserParams,
) (*user.User, error) {
	userModel, err := u.queries.GetUserByCredentials(ctx,
		sqlc.GetUserByCredentialsParams{
			Username:     params.Username,
			PasswordHash: params.PasswordHash,
		},
	)
	return userModelToUser(userModel), err
}

func (u *userRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*user.User, error) {
	userModel, err := u.queries.GetUserByUsername(ctx, username)
	return userModelToUser(userModel), err
}

func (u *userRepository) ListUsers(
	ctx context.Context,
	params user.ListUsersParams,
) ([]*user.User, int64, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return []*user.User{}, 0, err
	}
	defer tx.Rollback()
	qtx := u.queries.WithTx(tx)
	userModels, err := qtx.ListUsers(
		ctx,
		sqlc.ListUsersParams{
			Limit:  int64(params.Limit),
			Offset: int64(params.Offset),
		},
	)
	if err != nil {
		return []*user.User{}, 0, err
	}
	total, err := qtx.GetUserCount(ctx)
	if err != nil {
		return []*user.User{}, 0, err
	}
	var res []*user.User
	for _, v := range userModels {
		res = append(res, userModelToUser(v))
	}
	tx.Commit()
	return res, total, nil
}

func (u *userRepository) UpdateUser(
	ctx context.Context,
	params user.UpdateUserParams,
) (*user.User, error) {
	locked := int64(0)
	if params.Locked {
		locked = 1
	}
	userModel, err := u.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		PasswordHash: params.PasswordHash,
		Locked:       locked,
		ExpiredAt:    *params.ExpiredAt,
		ID:           params.ID,
	})
	return userModelToUser(userModel), err
}
