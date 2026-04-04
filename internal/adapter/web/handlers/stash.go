package handlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
)

// add update stash and delete stash handlers
// also add some terraform helper handlers to check before apply

// CreateStash create stash
//
//	@Summary	Create Stash
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.CreateStashRequest	true	"stash to create"
//	@Success	201		{object}	jsonutil.Message
//	@Failure	400		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes [post]
//	@Security	BearerAuth
func CreateStash(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		dto, err := jsonutil.Read[dto.CreateStashRequest](r.Body)
		if err != nil {
			return ports.BadRequestError(nil)
		}
		if err := stashService.CreateStash(ctx, dto); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Created)
	}
}

// DeleteStash delete stash by id
//
//	@Summary	Delete Stash
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Stash ID"
//	@Success	200	{object}	jsonutil.Message
//	@Failure	404	{object}	jsonutil.Message
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	403	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/stashes/{id} [delete]
//	@Security	BearerAuth
func DeleteStash(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		if err := stashService.DeleteStash(ctx, stashUUID); err != nil {
			return err
		}

		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

// GetStashByID get stash by id
//
//	@Summary	Get Stash By ID
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Stash ID"
//	@Success	200	{object}	dto.StashResponse
//	@Failure	404	{object}	jsonutil.Message
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	403	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/stashes/{id} [get]
//	@Security	BearerAuth
func GetStashByID(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		stash, err := stashService.GetStashByID(ctx, stashUUID)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, stash)
	}
}

// ListStashes list stashes that user maintains or is member of
//
//	@Summary	List User Stashes
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	dto.ListMyStashesResponse
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/stashes [get]
//	@Security	BearerAuth
func ListStashes(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		res, err := stashService.ListMyStashes(ctx)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, res)
	}
}

// LockStash lock certain stash to restrict any secrets access
//
//	@Summary	Lock Stash
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Stash ID"
//	@Success	200	{object}	jsonutil.Message
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	403	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/stashes/{id}/lock [post]
//	@Security	BearerAuth
func LockStash(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		if err := stashService.Lock(ctx, stashUUID); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

// UnlockStash unlock certain stash to make secrets available for maintainer or members
//
//	@Summary	Unlock Stash
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string							true	"Stash ID"
//	@Param		request	body		handlers.UnlockStash.request	true	"stash password"
//	@Success	200		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id}/unlock [post]
//	@Security	BearerAuth
func UnlockStash(stashService ports.StashService) apiFunc {
	type request struct {
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		req, err := jsonutil.Read[request](r.Body)
		if err != nil {
			return ports.BadRequestError(nil)
		}

		if err := stashService.Unlock(ctx, stashUUID, req.Password); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

// ListSecrets list secrets of an unlocked stash
//
//	@Summary	List Secrets
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Stash ID"
//	@Success	200	{object}	dto.SecretResponse
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	403	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/stashes/{id}/secrets [get]
//	@Security	BearerAuth
func GetSecrets(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		resp, err := stashService.GetSecrets(ctx, stashUUID)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, resp)
	}
}

// GetSecretsEntry get certain secrets entry from an unlocked stash
//
//	@Summary	Get Secrets Entry
//	@Tags		Stashes
//	@Accept		json
//	@Produce	plain
//	@Param		id		path		string	true	"Stash ID"
//	@Param		name	path		string	true	"Secrets entry name"
//	@Success	200		{object}	string
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id}/secrets/{name} [get]
//	@Security	BearerAuth
func GetSecretsEntry(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		entryName := r.PathValue("entry_name")

		resp, err := stashService.GetSecretsEntry(ctx, stashUUID, entryName)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))

		return nil
	}
}

// AddSecretsEntry add secrets entry
//
//	@Summary	Add Secrets Entry
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string			true	"Stash ID"
//	@Param		request	body		dto.AddSecret	true	"secret name and value"
//	@Success	201		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id}/secrets [put]
//	@Security	BearerAuth
func AddSecretsEntry(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		req := dto.AddSecret{}
		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			req, err = jsonutil.Read[dto.AddSecret](r.Body)
			if err != nil {
				return ports.BadRequestError(nil)
			}
		} else {
			if strings.HasPrefix(contentType, "multipart/form-data") {
				r.ParseMultipartForm(1 << 20)
			} else {
				r.ParseForm()
			}
			req = dto.AddSecret{
				Name:  r.FormValue("name"),
				Value: r.FormValue("value"),
			}
		}

		if err := stashService.AddSecretsEntry(ctx, stashUUID, req); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Created)
	}
}

// RemoveSecretsEntry remove secrets entry
//
//	@Summary	Remove Secrets Entry
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string	true	"Stash ID"
//	@Param		name	path		string	true	"Secrets entry name"
//	@Success	200		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id}/secrets/{name} [delete]
//	@Security	BearerAuth
func RemoveSecretsEntry(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		entryName := r.PathValue("entry_name")

		if err := stashService.RemoveSecretsEntry(ctx, stashUUID, entryName); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

// GetStashMembers get stash member
//
//	@Summary	Get stash members
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Stash ID"
//	@Success	200	{object}	dto.ListStashMemberResponse
//	@Failure	404	{object}	jsonutil.Message
//	@Failure	401	{object}	jsonutil.Message
//	@Failure	403	{object}	jsonutil.Message
//	@Failure	500	{object}	jsonutil.Message
//	@Router		/stashes/{id}/members [get]
//	@Security	BearerAuth
func GetStashMembers(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		resp, err := stashService.ListStashMembers(ctx, stashUUID)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, resp)
	}
}

// AddStashMember add stash member
//
//	@Summary	Add stash member
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"Stash ID"
//	@Param		request	body		dto.AddStashMemberRequest	true	"add stash member request"
//	@Success	201		{object}	jsonutil.Message
//	@Failure	404		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id}/members [post]
//	@Security	BearerAuth
func AddStashMember(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		req, err := jsonutil.Read[dto.AddStashMemberRequest](r.Body)
		if err != nil {
			return ports.BadRequestError(nil)
		}

		req.StashID = stashUUID

		if err := stashService.AddStashMember(ctx, req); err != nil {
			return err
		}

		return jsonutil.SendMessage(w, jsonutil.Created)
	}
}

// RemoveStashMember remove stash member
//
//	@Summary	Remove stash member
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string	true	"Stash ID"
//	@Param		user_id	path		string	true	"Stash member ID"
//	@Success	200		{object}	jsonutil.Message
//	@Failure	404		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id}/members/{user_id} [delete]
//	@Security	BearerAuth
func RemoveStashMember(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		userID := r.PathValue("user_id")
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return ports.NotFoundError(errUserNotFound)
		}

		req := dto.RemoveStashMemberRequest{
			StashID: stashUUID,
			UserID:  userUUID,
		}

		if err := stashService.RemoveStashMember(ctx, req); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

// UpdateStash update stash
//
//	@Summary	Update Stash
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"Stash ID"
//	@Param		request	body		dto.UpdateStashRequest	true	"update stash request"
//	@Success	200		{object}	jsonutil.Message
//	@Failure	404		{object}	jsonutil.Message
//	@Failure	401		{object}	jsonutil.Message
//	@Failure	403		{object}	jsonutil.Message
//	@Failure	500		{object}	jsonutil.Message
//	@Router		/stashes/{id} [patch]
//	@Security	BearerAuth
func UpdateStash(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(nil)
		}

		req, err := jsonutil.Read[dto.UpdateStashRequest](r.Body)
		if err != nil {
			return ports.BadRequestError(nil)
		}
		if err := stashService.UpdateStash(ctx, stashUUID, req); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

// GetStashByName get stash by name
//
//	@Summary	Get Stash By Name
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		maintainer_id	path		string	true	"Stash maintainer id"
//	@Param		name			path		string	true	"Stash name"
//	@Success	200				{object}	dto.StashResponse
//	@Failure	404				{object}	jsonutil.Message
//	@Failure	401				{object}	jsonutil.Message
//	@Failure	403				{object}	jsonutil.Message
//	@Failure	500				{object}	jsonutil.Message
//	@Router		/stashes/by-name/{maintainer_id}/{name} [get]
//	@Security	BearerAuth
func GetStashByName(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		maintainerID := r.PathValue("maintainer_id")
		maintainerUUID, err := uuid.Parse(maintainerID)
		if err != nil {
			return ports.NotFoundError(nil)
		}
		stashName := r.PathValue("stash_name")
		resp, err := stashService.GetStashByName(ctx, maintainerUUID, stashName)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, 200, resp)
	}
}

// GetStashMember get stash member
//
//	@Summary	Get stash member
//	@Tags		Stashes
//	@Accept		json
//	@Produce	json
//	@Param		stash_id	path		string	true	"Stash ID"
//	@Param		user_id		path		string	true	"User ID"
//	@Success	200			{object}	dto.StashMemberResponse
//	@Failure	404			{object}	jsonutil.Message
//	@Router		/stashes/{id}/members/{user_id} [get]
//	@Security	BearerAuth
func GetStashMember(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		stashID := r.PathValue("stash_id")
		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		userID := r.PathValue("user_id")
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return ports.NotFoundError(errStashNotFound)
		}

		resp, err := stashService.GetStashMember(ctx, stashUUID, userUUID)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, resp)
	}
}
