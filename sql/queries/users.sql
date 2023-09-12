-- name: AuthenticateUser :one
SELECT id, hashed_password
FROM users
WHERE email = ?;

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
