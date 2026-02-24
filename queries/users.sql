-- name: ListUsers :many
SELECT * FROM users ORDER BY username LIMIT ? OFFSET ?; 

-- name: GetUserById :one 
SELECT * FROM users WHERE id = ?;

-- name: GetUserByCredentials :one
SELECT * FROM users WHERE username = ? AND password_hash = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: AddUser :one
INSERT INTO users(username, password_hash)
VALUES (?, ?) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET password_hash = ?, locked = ?, expired_at = ? WHERE id = ? RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = ?;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users ORDER BY username;
