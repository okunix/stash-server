package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/stash"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type stashRepository struct {
	db *sql.DB
}

func NewStashRepository(db *sql.DB) ports.StashRepository {
	return &stashRepository{db: db}
}

func (s *stashRepository) AddMember(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

func (s *stashRepository) CommitData(ctx context.Context, params stash.CommitDataParams) error {
	panic("unimplemented")
}

func (s *stashRepository) CreateStash(
	ctx context.Context,
	params stash.CreateStashParams,
) (*stash.Stash, error) {
	panic("unimplemented")
}

func (s *stashRepository) DeleteStash(ctx context.Context, stashID uuid.UUID) error {
	panic("unimplemented")
}

func (s *stashRepository) GetStashByID(ctx context.Context, id uuid.UUID) (*stash.Stash, error) {
	panic("unimplemented")
}

func (s *stashRepository) ListStashes(
	ctx context.Context,
	params stash.ListStashesParams,
) ([]*stash.Stash, int64, error) {
	panic("unimplemented")
}

func (s *stashRepository) UpdateStash(
	ctx context.Context,
	params stash.UpdateStashParams,
) (*stash.Stash, error) {
	panic("unimplemented")
}
