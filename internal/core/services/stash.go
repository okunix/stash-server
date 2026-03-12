package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/dto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type stashService struct {
	stashRepo ports.StashRepository
}

type StashServiceParams struct {
	stashRepo ports.StashRepository
}

func NewStashService(params StashServiceParams) ports.StashService {
	return &stashService{stashRepo: params.stashRepo}
}

func (s *stashService) AddStashMember(ctx context.Context, req dto.AddStashMemberRequest) error {
	panic("unimplemented")
}

func (s *stashService) CheckStashMember(
	ctx context.Context,
	stashID uuid.UUID,
	userID uuid.UUID,
) (bool, error) {
	panic("unimplemented")
}

func (s *stashService) CreateStash(ctx context.Context, req dto.CreateStashRequest) error {
	panic("unimplemented")
}

func (s *stashService) DeleteStash(ctx context.Context, stashID uuid.UUID) error {
	panic("unimplemented")
}

func (s *stashService) GetStashByID(
	ctx context.Context,
	stashID uuid.UUID,
) (*dto.StashResponse, error) {
	panic("unimplemented")
}

func (s *stashService) ListStashMembers(
	ctx context.Context,
	stashID uuid.UUID,
) (*dto.ListStashMemberResponse, error) {
	panic("unimplemented")
}

func (s *stashService) ListStashes(
	ctx context.Context,
	req dto.ListStashesRequest,
) ([]*dto.StashResponse, error) {
	panic("unimplemented")
}

func (s *stashService) RemoveStashMember(
	ctx context.Context,
	req dto.RemoveStashMemberRequest,
) error {
	panic("unimplemented")
}

func (s *stashService) UpdateStash(ctx context.Context, req dto.UpdateStashRequest) error {
	panic("unimplemented")
}
