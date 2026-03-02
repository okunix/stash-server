package dto

import (
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
