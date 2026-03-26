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
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
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
