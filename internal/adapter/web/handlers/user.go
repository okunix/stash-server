package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/adapter/web/webutil"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
)

// GetUserByUsernameOrID get user by username or id
//
//	@Summary		Get username by username or id
//	@Description	Returns user by id(uuid) or (if request param is not uuid) by username
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			username_or_id	path		string	true	"Username or ID"
//	@Success		200				{object}	dto.UserResponse
//	@Failure		404				{object}	jsonutil.Message
//	@Failure		401				{object}	jsonutil.Message
//	@Failure		403				{object}	jsonutil.Message
//	@Failure		500				{object}	jsonutil.Message
//	@Router			/users/{username_or_id} [get]
//	@Security		BearerAuth
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

// Whoami get current user information
//
//	@Summary	Get current user information
//	@Tags		Accounts
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	dto.UserResponse
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/auth/whoami [get]
//	@Security	BearerAuth
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

// ChangePassword change password for current user
//
//	@Summary		Change Password
//	@Description	change password for current user or (if admin) for specified user
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ChangePasswordRequest	true	"change password request"
//	@Success		200		{object}	jsonutil.Message
//	@Failure		404		{object}	jsonutil.Message
//	@Failure		401		{object}	jsonutil.Message
//	@Failure		403		{object}	jsonutil.Message
//	@Failure		500		{object}	jsonutil.Message
//	@Router			/auth/change-password [patch]
//	@Security		BearerAuth
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

// ListUsers list users
//
//	@Summary	List users
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int	false	"limit"		minimum(0)
//	@Param		offset	query		int	false	"offset"	minimum(0)
//	@Success	200		{object}	dto.ListUsersResponse
//	@Failure	404		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/users [get]
//	@Security	BearerAuth
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
