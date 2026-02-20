-- name: ListAccessLog :many
SELECT sqlc.embed(u), sqlc.embed(s), a.secret_name, a.action, a.created_at FROM access_log a INNER JOIN users u ON u.id = user_id INNER JOIN stashes s ON s.id = stash_id ORDER BY a.created_at DESC LIMIT ? OFFSET ?;
