package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/auth"
	"gitlab.com/stash-password-manager/stash-server/internal/core/crypto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/user"
	"gitlab.com/stash-password-manager/stash-server/internal/core/dto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type userService struct {
	userRepo ports.UserRepository
}

type UserServiceParams struct {
	UserRepository ports.UserRepository
}

func NewUserService(params UserServiceParams) ports.UserService {
	return &userService{
		userRepo: params.UserRepository,
	}
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials provided")
)

func (u *userService) hashPassword(password string) (string, error) {
	hashFunc, err := crypto.NewArgon2ID()
	if err != nil {
		return "", err
	}
	hash, err := hashFunc.DeriveKey([]byte(password))
	if err != nil {
		return "", err
	}
	passwordHash := hash.String()
	return passwordHash, nil
}

func (u *userService) comparePasswordHash(hash, password string) (bool, error) {
	kdf, existingHash, err := crypto.NewArgon2IDFromString(hash)
	if err != nil {
		return false, err
	}
	passwordHash, err := kdf.DeriveKey([]byte(password))
	if err != nil {
		return false, err
	}
	return kdf.Compare(existingHash, passwordHash.Bytes()), nil
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) error {
	if problems, ok := req.Validate(); !ok {
		return ports.ValidationError(problems)
	}

	passwordHash, err := u.hashPassword(req.Password)
	if err != nil {
		slog.Error("password hash failed", "error", err.Error())
		return err
	}

	_, err = u.userRepo.AddUser(ctx, user.AddUserParams{
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
	return err
}

// use only in admin cli
func (u *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	slog.Info("deleting user", "id", userID)
	return u.DeleteUser(ctx, userID)
}

func (u *userService) GetUserToken(
	ctx context.Context,
	req dto.GetUserTokenRequest,
) (string, error) {
	slog.Info("retriving jwt token for user", "username", req.Username)

	user, err := u.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		slog.Warn("failed user lookup", "username", req.Username)
		return "", errors.New("invalid credentials")
	}

	ok, _ := u.comparePasswordHash(user.PasswordHash, req.Password)
	if !ok {
		return "", errors.New("invalid credentials")
	}

	slog.Info("generating jwt for user", "user_id", user.ID, "username", user.Username)
	return auth.JWT(user.ID, user.Username)
}

func (u *userService) LockUser(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

func (u *userService) UnlockUser(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}
