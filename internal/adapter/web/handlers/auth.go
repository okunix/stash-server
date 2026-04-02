package handlers

import (
	"net/http"

	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
)

// Login godoc
//
//	@Summary		Login User
//	@Description	Get User JWT Token for Login
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.GetUserTokenRequest	true	"user credentials"
//	@Success		200		{object}	dto.GetUserTokenResponse
//	@Failure		400		{object}	jsonutil.Message
//	@Failure		500		{object}	jsonutil.Message
//	@Router			/auth/login [post]
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

// CreateUser godoc
//
//	@Summary		Create User
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateUserRequest	true	"user credentials"
//	@Success		200		{object}	jsonutil.Message
//	@Failure		400		{object}	jsonutil.Message
//	@Failure		401		{object}	jsonutil.Message
//	@Failure		403		{object}	jsonutil.Message
//	@Failure		500		{object}	jsonutil.Message
//	@Router			/users [post]
//	@Security		BearerAuth
//	@Description	Requires admin role.
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
