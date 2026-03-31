package handlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/adapter/web/jsonutil"
	"github.com/okunix/stash-server/internal/adapter/web/webutil"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
)

// add update stash and delete stash handlers
// also add some terraform helper handlers to check before apply

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

func ListStashes(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		limit, _ := webutil.GetUintQueryParam(r, "limit", 32, 30)
		offset, _ := webutil.GetUintQueryParam(r, "offset", 32, 0)
		req := dto.ListStashesRequest{
			Limit:  uint(limit),
			Offset: uint(offset),
			Search: r.URL.Query().Get("search"),
		}
		res, err := stashService.ListStashes(ctx, req)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, res)
	}
}

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
