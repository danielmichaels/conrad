-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id              INTEGER     NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at      TEXT        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TEXT        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email           TEXT UNIQUE NOT NULL,
    hashed_password TEXT        NOT NULL,
    name            TEXT
);
CREATE TABLE gitlab_clients
(
    id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    created_at   TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by   INTEGER NOT NULL,
    webhook_url  TEXT    NOT NULL,
    gitlab_url   TEXT    NOT NULL,
    -- Time in seconds
    interval     INTEGER NOT NULL DEFAULT 86400,
    access_token TEXT    NOT NULL
);
CREATE TABLE gitlab_repos
(
    id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    repo_id      INTEGER NOT NULL,
    repo_web_url TEXT    NOT NULL,
    -- tracked: 0 is false, 1 is true
    tracked      INTEGER NOT NULL DEFAULT 0,
    name         TEXT    NOT NULL,
    created_at   TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id    INTEGER NOT NULL,
    FOREIGN KEY (client_id) REFERENCES gitlab_clients (id) ON DELETE CASCADE,
    UNIQUE (repo_id, client_id)
);
CREATE TRIGGER updated_at_users
    AFTER UPDATE
    ON users
    FOR EACH ROW
BEGIN
    UPDATE users
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TRIGGER updated_at_gitlab_clients
    AFTER UPDATE
    ON gitlab_clients
    FOR EACH ROW
BEGIN
    UPDATE gitlab_clients
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS updated_at_users;
DROP TRIGGER IF EXISTS updated_at_gitlab_clients;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS gitlab_clients;
DROP TABLE IF EXISTS gitlab_repos;
-- +goose StatementEnd
