package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/jsonutil"
	"gitlab.com/stash-password-manager/stash-server/internal/adapter/web/webutil"
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
		stashID := r.PathValue("stash_id")

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
