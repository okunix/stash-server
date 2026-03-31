package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/core/domain/user"
	"github.com/okunix/stash-server/internal/core/ports"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{db: db}
}

type userSQLModel struct {
	id           uuid.UUID
	username     string
	passwordHash string
	locked       bool
	role         string
	expiredAt    sql.NullTime
	createdAt    time.Time
}

func (u *userSQLModel) Domain() *user.User {
	if u == nil {
		return nil
	}
	var expiredAt *time.Time
	if u.expiredAt.Valid {
		expiredAt = &u.expiredAt.Time
	}
	return &user.User{
		ID:           u.id,
		Username:     u.username,
		PasswordHash: u.passwordHash,
		Locked:       u.locked,
		Role:         u.role,
		ExpiredAt:    expiredAt,
		CreatedAt:    u.createdAt,
	}
}

func scanUserSQLRow(row scannable) (*userSQLModel, error) {
	var res userSQLModel
	err := row.Scan(
		&res.id,
		&res.username,
		&res.passwordHash,
		&res.locked,
		&res.role,
		&res.expiredAt,
		&res.createdAt,
	)
	return &res, err
}

const addUserStmt = `
INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id, username, password_hash, locked, role, expired_at, created_at;
`

func (u *userRepository) AddUser(
	ctx context.Context,
	params user.AddUserParams,
) (*user.User, error) {
	row := u.db.QueryRowContext(ctx, addUserStmt, params.Username, params.PasswordHash, params.Role)
	userModel, err := scanUserSQLRow(row)
	return userModel.Domain(), err
}

const deleteUserStmt = `
DELETE FROM users WHERE id = $1;
`

func (u *userRepository) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	res, err := u.db.ExecContext(ctx, deleteUserStmt, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

const updateUserStmt = `
UPDATE users SET password_hash = $2, locked = $3, expired_at = $4 WHERE id = $1 RETURNING id, username, password_hash, locked, role, expired_at, created_at;
`

func (u *userRepository) UpdateUser(
	ctx context.Context,
	params user.UpdateUserParams,
) (*user.User, error) {
	userModel, err := scanUserSQLRow(u.db.QueryRowContext(
		ctx,
		updateUserStmt,
		params.ID,
		params.PasswordHash,
		params.Locked,
		params.ExpiredAt,
	))
	return userModel.Domain(), err
}

const getUserStmt = `SELECT id, username, password_hash, locked, role, expired_at, created_at FROM users WHERE username = $1 AND password_hash = $2 LIMIT 1;`

func (u *userRepository) GetUser(
	ctx context.Context,
	params user.GetUserParams,
) (*user.User, error) {
	userModel, err := scanUserSQLRow(
		u.db.QueryRowContext(ctx,
			getUserStmt,
			params.Username,
			params.PasswordHash,
		))
	return userModel.Domain(), err
}

const getUserByIDStmt = `SELECT id, username, password_hash, locked, role, expired_at, created_at FROM users WHERE id = $1 LIMIT 1;`

func (u *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	userModel, err := scanUserSQLRow(u.db.QueryRowContext(ctx, getUserByIDStmt, id))
	return userModel.Domain(), err
}

const getUserByUsernameStmt = `SELECT id, username, password_hash, locked, role, expired_at, created_at FROM users WHERE username = $1 LIMIT 1;`

func (u *userRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*user.User, error) {
	userModel, err := scanUserSQLRow(u.db.QueryRowContext(ctx, getUserByUsernameStmt, username))
	return userModel.Domain(), err
}

const (
	listUsersStmt     = `SELECT id, username, password_hash, locked, role, expired_at, created_at FROM users LIMIT $1 OFFSET $2;`
	getTotalUsersStmt = `SELECT COUNT(*) FROM users;`
)

func (u *userRepository) ListUsers(
	ctx context.Context,
	params user.ListUsersParams,
) ([]*user.User, int64, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return []*user.User{}, 0, err
	}
	defer tx.Rollback()

	users := make([]*user.User, 0, params.Limit)
	rows, err := tx.QueryContext(ctx, listUsersStmt, params.Limit, params.Offset)
	if err != nil {
		return users, 0, err
	}
	for rows.Next() {
		userModel, err := scanUserSQLRow(rows)
		if err != nil {
			slog.Info("failed to scan user sql model", "error", err.Error())
			continue
		}
		users = append(users, userModel.Domain())
	}

	var count int64
	if err := tx.QueryRowContext(ctx, getTotalUsersStmt).Scan(&count); err != nil {
		return users, 0, err
	}

	tx.Commit()
	return users, count, nil
}

const isAdminPresentSQL = `
	SELECT EXISTS (
		SELECT 1 FROM users WHERE role = 'admin'
	);
`

func (u *userRepository) IsAdminPresent(ctx context.Context) (bool, error) {
	var exists bool
	err := u.db.QueryRowContext(ctx, isAdminPresentSQL).Scan(&exists)
	return exists, err
}
