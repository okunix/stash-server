package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/core/dto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

func CreateStash(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		dto, err := jsonutil.Read[dto.CreateStashRequest](r.Body)
		if err != nil {
			return ports.BadRequestError(err)
		}
		if err := stashService.CreateStash(ctx, dto); err != nil {
			return err
		}
		return jsonutil.SendMessage(w, jsonutil.Ok)
	}
}

func DeleteStash(stashService ports.StashService) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		stashID := r.PathValue("id")

		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.BadRequestError(err)
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
		stashID := r.PathValue("id")

		stashUUID, err := uuid.Parse(stashID)
		if err != nil {
			return ports.BadRequestError(err)
		}

		stash, err := stashService.GetStashByID(ctx, stashUUID)
		if err != nil {
			return err
		}
		return jsonutil.Write(w, http.StatusOK, stash)
	}
}
