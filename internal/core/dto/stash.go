package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/core/domain/stash"
)

type StashResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	MaintainerID uuid.UUID `json:"maintainer_id"`
	CreatedAt    time.Time `json:"created_at"`
	Locked       bool      `json:"locked"`
}

func NewStashResponse(s *stash.Stash, locked bool) *StashResponse {
	return &StashResponse{
		ID:           s.ID,
		Name:         s.Name,
		Description:  s.Description,
		MaintainerID: s.MaintainerID,
		CreatedAt:    s.CreatedAt,
		Locked:       locked,
	}
}

type ListStashResponse struct {
	Page   *Page           `json:"page,omitempty"`
	Result []StashResponse `json:"result"`
}

type GetStashByIDRequest struct {
	StashID uuid.UUID `json:"stash_id"`
}

type CreateStashRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Password    string  `json:"password"`
}

func (req CreateStashRequest) Validate() (map[string]string, bool) {
	problems := make(map[string]string)
	if err := stash.ValidateName(req.Name); err != nil {
		problems["name"] = err.Error()
	}
	if err := stash.ValidateDescription(req.Description); err != nil {
		problems["description"] = err.Error()
	}
	if err := stash.ValidatePassword(req.Password); err != nil {
		problems["password"] = err.Error()
	}
	return problems, len(problems) == 0
}

type UpdateStashRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (req UpdateStashRequest) Validate() (map[string]string, bool) {
	problems := make(map[string]string)
	if req.Name != nil {
		if err := stash.ValidateName(*req.Name); err != nil {
			problems["name"] = err.Error()
		}
	}
	if err := stash.ValidateDescription(req.Description); err != nil {
		problems["description"] = err.Error()
	}
	return problems, len(problems) == 0
}

type ListStashesRequest struct {
	Limit  uint   `json:"limit"`
	Offset uint   `json:"offset"`
	Search string `json:"search"`
}

type StashMemberResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Since    time.Time `json:"since"`
}

type StashMaintainerResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

type ListStashMemberResponse struct {
	Maintainer StashMaintainerResponse `json:"maintainer"`
	Members    []StashMemberResponse   `json:"members"`
}

type AddStashMemberRequest struct {
	StashID uuid.UUID `json:"stash_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type RemoveStashMemberRequest struct {
	StashID uuid.UUID `json:"stash_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type ListMyStashesResponse struct {
	Maintainer []*StashResponse `json:"maintainer"`
	Member     []*StashResponse `json:"member"`
}
