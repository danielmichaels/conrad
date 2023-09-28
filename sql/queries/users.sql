-- name: AuthenticateUser :one
SELECT id, hashed_password
FROM users
-- it's always 1 because we only have a single user in Conrad v1
WHERE id = 1;

-- name: InitialisePassphrase :one
INSERT OR REPLACE INTO users
    (id, hashed_password)
VALUES (?, ?)
RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: InsertNewUser :one
INSERT INTO users
    (name, email, hashed_password)
VALUES (?, ?, ?)
RETURNING id;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = ?;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;
