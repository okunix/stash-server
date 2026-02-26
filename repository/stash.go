package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/domain/stash"
	"gitlab.com/stash-password-manager/stash-server/domain/user"
	"gitlab.com/stash-password-manager/stash-server/sqlc"
)

type stashRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

func NewStashRepository(db *sql.DB) stash.Repository {
	return &stashRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func stashModelToStash(stashModel *sqlc.Stash, maintainerModel *sqlc.User) *stash.Stash {
	if stashModel == nil {
		return nil
	}
	return &stash.Stash{
		ID:            stashModel.ID,
		Name:          stashModel.Name,
		MaintainerID:  stashModel.MaintainerID,
		Description:   stashModel.Description,
		MasterKeyHash: stashModel.MasterKeyHash,
		EncryptedData: *stashModel.EncryptedData,
		CreatedAt:     stashModel.CreatedAt,
	}
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

func (s *stashRepository) GetStash(
	ctx context.Context,
	stashID uuid.UUID,
) (*stash.Stash, error) {
	panic("unimplemented")
}

func (s *stashRepository) ListStashes(
	ctx context.Context,
	params stash.ListStashesParams,
) ([]*stash.Stash, int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return []*stash.Stash{}, 0, err
	}
	defer tx.Rollback()
	queries := s.queries.WithTx(tx)
	rows, err := queries.ListStashes(
		ctx,
		sqlc.ListStashesParams{
			MaintainerID: params.MaintainerID,
			Limit:        int64(params.Limit),
			Offset:       int64(params.Offset),
		},
	)
	if err != nil || rows == nil {
		return []*stash.Stash{}, 0, err
	}
	count, err := queries.GetStashesCount(ctx, params.MaintainerID)
	if err != nil {
		return []*stash.Stash{}, 0, err
	}
	var stashes []*stash.Stash
	for _, v := range rows {
		stash := stashModelToStash(&v.Stash, &v.User)
		userModels, err := queries.ListStashMembers(ctx, stash.ID)
		if err != nil {
			slog.Error("failed to get some users", "error", err.Error())
			continue
		}
		members := []user.User{}
		for _, user := range userModels {
			member := userModelToUser(user)
			members = append(members, *member)
		}
		stashes = append(stashes, stash)
	}
	tx.Commit()
	return stashes, count, nil
}

func (s *stashRepository) UpdateStash(
	ctx context.Context,
	params stash.UpdateStashParams,
) (*stash.Stash, error) {
	panic("unimplemented")
}
