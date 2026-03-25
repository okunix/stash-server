-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    locked BOOLEAN NOT NULL DEFAULT false,
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    expired_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS stashes (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    name TEXT NOT NULL,
    description TEXT,
    maintainer_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    master_key_hash TEXT NOT NULL,
    master_key_salt TEXT NOT NULL,
    encrypted_data TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    UNIQUE (name, maintainer_id)
);

CREATE TABLE IF NOT EXISTS stash_member (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stash_id UUID NOT NULL REFERENCES stashes(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (user_id, stash_id)
);

CREATE TABLE IF NOT EXISTS access_log (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    stash_id UUID NOT NULL REFERENCES stashes(id) ON DELETE NO ACTION,
    secret_name TEXT NOT NULL,
    action TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_log;
DROP TABLE IF EXISTS stash_member;
DROP TABLE IF EXISTS stashes;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
