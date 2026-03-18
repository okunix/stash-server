package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/accesslog"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/secret"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/stash"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/user"
)

type AccessLogRepository interface {
	ListLogs(
		ctx context.Context,
		params accesslog.ListLogsParams,
	) ([]*accesslog.AccessLog, int64, error)

	RemoveLogs(
		ctx context.Context,
		timestamp *time.Time,
	) (int64, error)

	AddLog(
		ctx context.Context,
		params accesslog.CreateLogEntryParams,
	) error
}

type SecretRepository interface {
	AddSecret(ctx context.Context, params secret.AddSecretParams) (*secret.Secret, error)
	RemoveSecretByStashID(ctx context.Context, stashID uuid.UUID) (*secret.Secret, error)
	GetSecretByStashID(ctx context.Context, stashID uuid.UUID) (*secret.Secret, error)
	ListSecrets(ctx context.Context, maintainerID uuid.UUID) ([]*secret.Secret, error)
	UpdateSecret(ctx context.Context, stashID uuid.UUID, sec *secret.Secret) error
}

type StashRepository interface {
	ListStashes(ctx context.Context, params stash.ListStashesParams) ([]*stash.Stash, int64, error)
	GetStashByID(ctx context.Context, id uuid.UUID) (*stash.Stash, error)
	CreateStash(ctx context.Context, params stash.CreateStashParams) (*stash.Stash, error)
	UpdateStash(ctx context.Context, params stash.UpdateStashParams) (*stash.Stash, error)
	DeleteStash(ctx context.Context, stashID uuid.UUID) error
	AddMember(ctx context.Context, params stash.AddMemberParams) error
	RemoveMember(ctx context.Context, params stash.RemoveMemberParams) error
	CommitData(ctx context.Context, params stash.CommitDataParams) error
	GetStashMembers(ctx context.Context, stashID uuid.UUID) ([]*stash.StashMember, error)
	IsStashMember(ctx context.Context, userID, stashID uuid.UUID) (bool, error)
	IsStashMaintainer(ctx context.Context, userID, stashID uuid.UUID) (bool, error)
	IsStashMemberOrMaintainer(ctx context.Context, userID, stashID uuid.UUID) (bool, error)
}

type UserRepository interface {
	ListUsers(ctx context.Context, params user.ListUsersParams) ([]*user.User, int64, error)
	GetUser(ctx context.Context, params user.GetUserParams) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	AddUser(ctx context.Context, params user.AddUserParams) (*user.User, error)
	UpdateUser(ctx context.Context, params user.UpdateUserParams) (*user.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
