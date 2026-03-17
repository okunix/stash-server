package ports

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/dto"
)

// current user information is transfered using context

type StashService interface {
	ListStashes(ctx context.Context, req dto.ListStashesRequest) (*dto.ListStashResponse, error)
	GetStashByID(ctx context.Context, stashID uuid.UUID) (*dto.StashResponse, error)
	CreateStash(ctx context.Context, req dto.CreateStashRequest) error
	UpdateStash(ctx context.Context, req dto.UpdateStashRequest) error
	DeleteStash(ctx context.Context, stashID uuid.UUID) error

	ListStashMembers(ctx context.Context, stashID uuid.UUID) (*dto.ListStashMemberResponse, error)
	AddStashMember(ctx context.Context, req dto.AddStashMemberRequest) error
	RemoveStashMember(ctx context.Context, req dto.RemoveStashMemberRequest) error

	GetSecrets(ctx context.Context, stashID uuid.UUID) (*dto.SecretResponse, error)
	GetSecretsEntry(ctx context.Context, stashID uuid.UUID, entryKey string) (string, error)
	ListUnlockedSecrets(ctx context.Context) ([]*dto.SecretResponse, error)

	Unlock(ctx context.Context, stashID uuid.UUID, password string) error
	Lock(ctx context.Context, stashID uuid.UUID) error
}

type UserService interface {
	GetUserToken(ctx context.Context, req dto.GetUserTokenRequest) (string, error)

	// admin cli functions
	CreateUser(ctx context.Context, req dto.CreateUserRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	LockUser(ctx context.Context, userID uuid.UUID) error
	UnlockUser(ctx context.Context, userID uuid.UUID) error
}

type AccessLogService interface {
	AddLogEntry(ctx context.Context, req dto.AddLogEntryRequest) error
	ListLogs(ctx context.Context, req dto.ListAccessLogRequest) error
}
