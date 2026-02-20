-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id TEXT PRIMARY KEY DEFAULT (generate_uuid()),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    password_salt TEXT NOT NULL,
    locked INTEGER NOT NULL CHECK (locked IN (0, 1)) DEFAULT 0,
    expired_at TEXT NULL, -- iso8601 
    created_at TEXT NOT NULL DEFAULT (datetime()) -- iso8601
);

CREATE TABLE IF NOT EXISTS stashes (
    id TEXT PRIMARY KEY DEFAULT (generate_uuid()),
    name TEXT NOT NULL UNIQUE,
    description TEXT NULL,
    maintainer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    master_key_hash TEXT NOT NULL,
    --master_key_salt TEXT NOT NULL UNIQUE,
    encrypted_data TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime()) -- iso8601
);

CREATE TABLE IF NOT EXISTS stash_member (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stash_id TEXT NOT NULL REFERENCES stashes(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (datetime()), -- iso8601
    PRIMARY KEY (user_id, stash_id)
);

CREATE TABLE IF NOT EXISTS access_log (
    id TEXT PRIMARY KEY DEFAULT (generate_uuid()),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    stash_id TEXT NOT NULL REFERENCES stashes(id) ON DELETE NO ACTION,
    secret_name TEXT NOT NULL,
    action TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime()) -- iso8601
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_log;
DROP TABLE IF EXISTS stash_member;
DROP TABLE IF EXISTS stashes;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
