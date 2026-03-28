package dto

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/user"
)

type GetUserTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req CreateUserRequest) Validate() (problems map[string]string, ok bool) {
	problems = make(map[string]string)
	if err := user.ValidateUsername(req.Username); err != nil {
		problems["username"] = err.Error()
	}
	if err := user.ValidatePassword(req.Password); err != nil {
		problems["password"] = err.Error()
	}
	return problems, len(problems) == 0
}

type GetUserTokenResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Locked    bool       `json:"locked"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}

type InitAdminResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	OldPassword string     `json:"old_password"`
	NewPassword string     `json:"new_password"`
}

func (req ChangePasswordRequest) Validate() (map[string]string, bool) {
	problems := make(map[string]string)
	if err := user.ValidatePassword(req.NewPassword); err != nil {
		problems["new_password"] = err.Error()
	}
	return problems, len(problems) == 0
}

func NewUserResponse(d *user.User) *UserResponse {
	if d == nil {
		return nil
	}
	return &UserResponse{
		ID:        d.ID,
		Username:  d.Username,
		Locked:    d.Locked,
		CreatedAt: d.CreatedAt,
		ExpiredAt: d.ExpiredAt,
	}
}

type GetUsersResponse struct {
	Page    *Page           `json:"page,omitempty"`
	Content []*UserResponse `json:"content"`
}

type GetUsersRequest struct {
	Limit  uint
	Offset uint
}
