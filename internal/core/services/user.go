package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/core/auth"
	"github.com/okunix/stash-server/internal/core/crypto"
	"github.com/okunix/stash-server/internal/core/domain/user"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
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

func (u *userService) checkAdminUser(ctx context.Context) error {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return ports.UnauthorizedError(nil)
	}
	if currentUser.Role != user.RoleAdmin {
		return ports.ForbiddenError(nil)
	}
	return nil
}

func (u *userService) createUserWithRole(
	ctx context.Context,
	req dto.CreateUserRequest,
	role string,
) error {
	if err := u.checkAdminUser(ctx); err != nil {
		return err
	}

	if problems, ok := req.Validate(); !ok {
		return ports.NewValidationError(problems)
	}

	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		slog.Error("password hash failed", "error", err.Error())
		return ports.InternalError(err)
	}

	_, err = u.userRepo.AddUser(ctx, user.AddUserParams{
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         role,
	})
	if err != nil {
		return ports.BadRequestError(err)
	}
	return nil
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) error {
	return u.createUserWithRole(ctx, req, user.RoleUser)
}

func (u *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := u.checkAdminUser(ctx); err != nil {
		return err
	}

	slog.Info("deleting user", "id", userID)
	if err := u.userRepo.DeleteUser(ctx, userID); err != nil {
		return ports.NotFoundError(err)
	}
	return nil
}

func (u *userService) GetUserToken(
	ctx context.Context,
	req dto.GetUserTokenRequest,
) (string, error) {
	slog.Info("retriving jwt token", "username", req.Username)

	user, err := u.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		slog.Warn("failed user lookup", "username", req.Username)
		return "", ports.BadRequestError(errors.New("invalid credentials"))
	}

	ok, _ := crypto.ComparePasswordHash(user.PasswordHash, req.Password)
	if !ok {
		return "", ports.BadRequestError(errors.New("invalid credentials"))
	}

	slog.Info("generating jwt", "user_id", user.ID, "username", user.Username)
	expiresAt := time.Now().Add(5 * time.Minute)
	token, err := auth.JWT(user.ID, user.Username, user.Role, auth.WithExpirationTime(expiresAt))
	if err != nil {
		return "", ports.InternalError(err)
	}
	return token, nil
}

func (u *userService) LockUser(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

func (u *userService) UnlockUser(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

func (u *userService) GetCurrentUser(ctx context.Context) (*dto.UserResponse, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, ports.UnauthorizedError(nil)
	}
	return u.GetUserByID(ctx, currentUser.UserID)
}

func (u *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	usr, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponse(usr), nil
}

func (u *userService) GetUserByUsername(
	ctx context.Context,
	username string,
) (*dto.UserResponse, error) {
	usr, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, ports.NotFoundError(errors.New("user not found"))
	}
	return dto.NewUserResponse(usr), nil
}

func (u *userService) InitializeAdminUser(ctx context.Context) (*dto.InitAdminResponse, error) {
	ok, err := u.userRepo.IsAdminPresent(ctx)
	if err != nil {
		return nil, ports.InternalError(err)
	}
	if ok {
		return nil, nil
	}
	initialAdminPasswordLength := 23
	initialAdminPassword := crypto.RandomSpecialString(initialAdminPasswordLength)
	req := dto.CreateUserRequest{
		Username: "admin_" + crypto.RandomAlphaNumericString(5),
		Password: initialAdminPassword,
	}
	err = u.createUserWithRole(ctx, req, user.RoleAdmin)
	return &dto.InitAdminResponse{
		Username: req.Username,
		Password: req.Password,
	}, err
}

func (u *userService) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return ports.UnauthorizedError(nil)
	}

	if problems, ok := req.Validate(); !ok {
		return ports.NewValidationError(problems)
	}

	userID := currentUser.UserID
	if req.UserID != nil {
		if currentUser.Role != user.RoleAdmin {
			return ports.ForbiddenError(nil)
		}
		userID = *req.UserID
	}
	userFromRepo, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return ports.NotFoundError(err)
	}

	ok, err = crypto.ComparePasswordHash(userFromRepo.PasswordHash, req.OldPassword)
	if err != nil {
		return ports.InternalError(err)
	}
	if !ok {
		return ports.BadRequestError(errors.New("old password is incorrect"))
	}

	newPasswordHash, _ := crypto.HashPassword(req.NewPassword)
	_, err = u.userRepo.UpdateUser(
		ctx,
		user.UpdateUserParams{
			ID:           userFromRepo.ID,
			PasswordHash: newPasswordHash,
			Locked:       userFromRepo.Locked,
			ExpiredAt:    userFromRepo.ExpiredAt,
		},
	)
	if err != nil {
		return ports.InternalError(err)
	}
	return nil
}

func (u *userService) ListUsers(
	ctx context.Context,
	req dto.ListUsersRequest,
) (*dto.ListUsersResponse, error) {
	_, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, ports.UnauthorizedError(nil)
	}

	users, total, err := u.userRepo.ListUsers(ctx,
		user.ListUsersParams{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
	if err != nil {
		return nil, ports.InternalError(err)
	}

	resp := dto.ListUsersResponse{
		Page: &dto.Page{
			Limit:  req.Limit,
			Offset: req.Offset,
			Total:  total,
		},
		Result: []*dto.UserResponse{},
	}

	for _, u := range users {
		resp.Result = append(resp.Result, dto.NewUserResponse(u))
	}

	return &resp, nil
}
