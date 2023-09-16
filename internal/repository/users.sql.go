// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: users.sql

package repository

import (
	"context"
	"database/sql"
)

const authenticateUser = `-- name: AuthenticateUser :one
SELECT id, hashed_password
FROM users
WHERE email = ?
`

type AuthenticateUserRow struct {
	ID             int64  `json:"id"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) AuthenticateUser(ctx context.Context, email string) (AuthenticateUserRow, error) {
	row := q.db.QueryRowContext(ctx, authenticateUser, email)
	var i AuthenticateUserRow
	err := row.Scan(&i.ID, &i.HashedPassword)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, name
FROM users
WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.Name,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, updated_at, email, hashed_password, name
FROM users
WHERE id = ?
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.Name,
	)
	return i, err
}

const insertNewUser = `-- name: InsertNewUser :one
INSERT INTO users
    (name, email, hashed_password)
VALUES (?, ?, ?)
RETURNING id
`

type InsertNewUserParams struct {
	Name           sql.NullString `json:"name"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"hashed_password"`
}

func (q *Queries) InsertNewUser(ctx context.Context, arg InsertNewUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertNewUser, arg.Name, arg.Email, arg.HashedPassword)
	var id int64
	err := row.Scan(&id)
	return id, err
}