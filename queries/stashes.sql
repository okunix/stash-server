-- name: ListStashes :many
SELECT sqlc.embed(stashes), sqlc.embed(users) FROM stashes INNER JOIN users ON users.id = stashes.maintainer_id WHERE stashes.maintainer_id = ? LIMIT ? OFFSET ?;

-- name: ListStashMembers :many
SELECT u.* FROM stash_member s INNER JOIN users u ON m.id = s.user_id WHERE stash_id = ?;

-- name: GetStashByID :one
SELECT sqlc.embed(stashes), sqlc.embed(users) FROM stashes INNER JOIN users ON stashes.maintainer_id = users.id WHERE stashes.id = ? ;

-- name: UpdateStash :execrows
UPDATE stashes SET name = ?, description = ?, master_key_hash = ?, encrypted_data = ? WHERE id = ?;

-- name: UpdateEncryptedData :execrows
UPDATE stashes SET encrypted_data = ? WHERE id = ?;

-- name: CreateStash :one
INSERT INTO stashes (name, description, maintainer_id, master_key_hash, encrypted_data) VALUES (?, ?, ?, ?, ?) RETURNING id;

-- name: DeleteStash :execrows
DELETE FROM stashes WHERE id = ?;

-- name: GetStashesCount :one
SELECT COUNT(*) FROM stashes WHERE stashes.maintainer_id = ?;
