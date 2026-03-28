package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

func GetUserByID(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		userID := r.PathValue("user_id")
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return ports.NotFoundError(nil)
		}
		resp, err := userService.GetUserByID(ctx, userUUID)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, resp)
	}
}

func Whoami(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		resp, err := userService.GetCurrentUser(ctx)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, resp)
	}
}
