package handlers

import (
	"net/http"

	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
)

// add get current user(whoami), update password and delete user handlers

func Login(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		req, err := jsonutil.Read[dto.GetUserTokenRequest](r.Body)
		if err != nil {
			return err
		}
		token, err := userService.GetUserToken(ctx, req)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, dto.GetUserTokenResponse{Token: token})
	}
}

func CreateUser(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		req, err := jsonutil.Read[dto.CreateUserRequest](r.Body)
		if err != nil {
			return err
		}
		if err := userService.CreateUser(ctx, req); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}
