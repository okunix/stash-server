package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/stash"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type stashSQLModel struct {
	id            uuid.UUID
	name          string
	description   sql.NullString
	maintainerID  uuid.UUID
	masterKeyHash string
	encryptedData sql.NullString
	createdAt     time.Time
}

func (s *stashSQLModel) Domain() *stash.Stash {
	if s == nil {
		return nil
	}
	var desc *string
	if s.description.Valid {
		desc = &s.description.String
	}
	encryptedData := s.encryptedData.String
	return &stash.Stash{
		ID:            s.id,
		Name:          s.name,
		Description:   desc,
		MaintainerID:  s.maintainerID,
		MasterKeyHash: s.masterKeyHash,
		EncryptedData: encryptedData,
		CreatedAt:     s.createdAt,
	}
}

func scanStashSQLRow(row scannable) (*stashSQLModel, error) {
	var resp stashSQLModel
	err := row.Scan(
		&resp.id,
		&resp.name,
		&resp.description,
		&resp.maintainerID,
		&resp.masterKeyHash,
		&resp.encryptedData,
		&resp.createdAt,
	)
	return &resp, err
}

type stashRepository struct {
	db *sql.DB
}

func NewStashRepository(db *sql.DB) ports.StashRepository {
	return &stashRepository{db: db}
}

const addMemberStmt = `
INSERT INTO stash_member (user_id, stash_id) VALUES ($1, $2);
`

func (s *stashRepository) AddMember(ctx context.Context, params stash.AddMemberParams) error {
	_, err := s.db.ExecContext(ctx, addMemberStmt, params.UserID, params.StashID)
	return err
}

const commitDataStmt = `
UPDATE stashes SET encrypted_data = $1 WHERE id = $2;
`

func (s *stashRepository) CommitData(ctx context.Context, params stash.CommitDataParams) error {
	_, err := s.db.ExecContext(ctx, commitDataStmt, params.Data, params.StashID)
	return err
}

const createStashStmt = `
INSERT INTO stashes (name, description, maintainer_id, master_key_hash, encrypted_data)
VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;
`

func (s *stashRepository) CreateStash(
	ctx context.Context,
	params stash.CreateStashParams,
) (*stash.Stash, error) {
	var id uuid.UUID
	var createdAt time.Time
	err := s.db.QueryRowContext(
		ctx,
		createStashStmt,
		params.Name,
		params.Description,
		params.MaintainerID,
		params.MasterKeyHash,
		params.EncryptedData,
	).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	res := &stash.Stash{
		ID:            id,
		Name:          params.Name,
		Description:   params.Description,
		MaintainerID:  params.MaintainerID,
		MasterKeyHash: params.MasterKeyHash,
		EncryptedData: params.EncryptedData,
		CreatedAt:     createdAt,
	}
	return res, err
}

const deleteStashStmt = `
DELETE FROM stashes WHERE id = $1;
`

func (s *stashRepository) DeleteStash(ctx context.Context, stashID uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, deleteStashStmt, stashID)
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected <= 0 {
		return errors.New("stash not found")
	}
	return err
}

const updateStashStmt = `
UPDATE stashes SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description, maintainer_id, master_key_hash, encrypted_data, created_at;
`

func (s *stashRepository) UpdateStash(
	ctx context.Context,
	params stash.UpdateStashParams,
) (*stash.Stash, error) {
	stashUpdateSQLResp, err := scanStashSQLRow(
		s.db.QueryRowContext(ctx, updateStashStmt, params.Name, params.Description),
	)
	return stashUpdateSQLResp.Domain(), err

}

const getStashByIDStmt = `
SELECT id, name, description, maintainer_id, master_key_hash, encrypted_data, created_at FROM stashes WHERE id = $1;
`

func (s *stashRepository) GetStashByID(ctx context.Context, id uuid.UUID) (*stash.Stash, error) {
	stashSQLResp, err := scanStashSQLRow(s.db.QueryRowContext(ctx, getStashByIDStmt, id))
	return stashSQLResp.Domain(), err
}

const (
	getTotalStashesStmt = `SELECT COUNT(*) FROM stashes WHERE maintainer_id = $1;`
	listStashesStmt     = `SELECT id, name, description, maintainer_id, master_key_hash, encrypted_data, created_at FROM stashes WHERE maintainer_id = $1 LIMIT $2 OFFSET $3;`
)

func (s *stashRepository) ListStashes(
	ctx context.Context,
	params stash.ListStashesParams,
) ([]*stash.Stash, int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return []*stash.Stash{}, 0, err
	}
	defer tx.Rollback()

	var count int64
	err = tx.QueryRowContext(ctx, getTotalStashesStmt, params.MaintainerID).Scan(&count)
	if err != nil {
		return []*stash.Stash{}, 0, err
	}

	rows, err := tx.QueryContext(
		ctx,
		listStashesStmt,
		params.MaintainerID,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		return []*stash.Stash{}, 0, err
	}

	stashes := make([]*stash.Stash, 0, params.Limit)
	for rows.Next() {
		stashSQLResp, err := scanStashSQLRow(rows)
		if err != nil {
			slog.Error("failed to scan stash row", "error", err.Error())
			continue
		}
		stashes = append(stashes, stashSQLResp.Domain())
	}

	tx.Commit()
	return stashes, count, nil
}
