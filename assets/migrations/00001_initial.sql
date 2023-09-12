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
CREATE TRIGGER updated_at_users
    AFTER UPDATE
    ON users
    FOR EACH ROW
BEGIN
    UPDATE users
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS updated_at_users;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
