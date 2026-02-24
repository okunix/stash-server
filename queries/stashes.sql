-- name: ListStashes :many
SELECT * FROM stashes WHERE maintainer_id = ? LIMIT ? OFFSET ?;

-- name: GetStashByID :one
SELECT * FROM stashes WHERE id = ?;

-- name: UpdateStash :execrows
UPDATE stashes SET name = ?, description = ?, master_key_hash = ?, encrypted_data = ? WHERE id = ?;

-- name: UpdateEncryptedData :execrows
UPDATE stashes SET encrypted_data = ? WHERE id = ?;

-- name: CreateStash :one
INSERT INTO stashes (name, description, maintainer_id, master_key_hash, encrypted_data) VALUES (?, ?, ?, ?, ?) RETURNING id;

-- name: DeleteStash :execrows
DELETE FROM stashes WHERE id = ?;
