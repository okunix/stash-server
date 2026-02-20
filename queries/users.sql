-- name: ListUsers :many
SELECT * FROM users ORDER BY username; 

-- name: GetUserById :one 
SELECT * FROM users WHERE id = ?;

-- name: AddUser :one
INSERT INTO users(username, password_hash, password_salt)
VALUES (?, ?, ?) RETURNING *;

-- name: UpdateUser :execrows
UPDATE users SET username = ?, password_hash = ?, password_salt = ?, locked = ?, expired_at = ? WHERE id = ?;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = ?;

