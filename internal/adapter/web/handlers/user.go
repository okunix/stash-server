package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/webutil"
	"gitlab.com/stash-password-manager/stash-server/internal/core/dto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

func GetUserByUsernameOrID(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		userID := r.PathValue("user_id")
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			resp, err := userService.GetUserByUsername(ctx, userID)
			if err != nil {
				return err
			}
			return jsonutil.Write(w, http.StatusOK, resp)
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

func ChangePassword(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		req, err := jsonutil.Read[dto.ChangePasswordRequest](r.Body)
		if err != nil {
			return ports.BadRequestError(nil)
		}
		if err := userService.ChangePassword(ctx, req); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

func ListUsers(userService ports.UserService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		limit, _ := webutil.GetUintQueryParam(r, "limit", 32, 50)
		offset, _ := webutil.GetUintQueryParam(r, "offset", 32, 0)
		req := dto.ListUsersRequest{Limit: uint(limit), Offset: uint(offset)}
		resp, err := userService.ListUsers(ctx, req)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, resp)
	}
}
