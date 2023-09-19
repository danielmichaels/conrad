// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: clients.sql

package repository

import (
	"context"
)

const getAllClientRepos = `-- name: GetAllClientRepos :many
SELECT gr.repo_id, gr.repo_web_url, gr.tracked, gr.name, gr.created_at, gr.updated_at, gr.client_id
FROM gitlab_repos gr
         INNER JOIN gitlab_clients gc ON gr.client_id = gc.id
WHERE gc.id = ?
`

func (q *Queries) GetAllClientRepos(ctx context.Context, id int64) ([]GitlabRepos, error) {
	rows, err := q.db.QueryContext(ctx, getAllClientRepos, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GitlabRepos{}
	for rows.Next() {
		var i GitlabRepos
		if err := rows.Scan(
			&i.RepoID,
			&i.RepoWebUrl,
			&i.Tracked,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ClientID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllClients = `-- name: GetAllClients :many
SELECT gc.id, gc.name, gc.created_at, gc.updated_at, gc.created_by, gc.webhook_url, gc.gitlab_url, gc.insecure, gc.interval, gc.access_token, COUNT(gr.repo_id) AS repo_count
FROM gitlab_clients gc
         LEFT JOIN gitlab_repos gr ON gc.id = gr.client_id
GROUP BY gc.id
`

type GetAllClientsRow struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedBy   int64  `json:"created_by"`
	WebhookUrl  string `json:"webhook_url"`
	GitlabUrl   string `json:"gitlab_url"`
	Insecure    string `json:"insecure"`
	Interval    int64  `json:"interval"`
	AccessToken string `json:"access_token"`
	RepoCount   int64  `json:"repo_count"`
}

func (q *Queries) GetAllClients(ctx context.Context) ([]GetAllClientsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllClients)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllClientsRow{}
	for rows.Next() {
		var i GetAllClientsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.WebhookUrl,
			&i.GitlabUrl,
			&i.Insecure,
			&i.Interval,
			&i.AccessToken,
			&i.RepoCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getClientById = `-- name: GetClientById :one
SELECT id, name, created_at, updated_at, created_by, webhook_url, gitlab_url, insecure, interval, access_token
FROM gitlab_clients
WHERE id = ?
`

func (q *Queries) GetClientById(ctx context.Context, id int64) (GitlabClients, error) {
	row := q.db.QueryRowContext(ctx, getClientById, id)
	var i GitlabClients
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.WebhookUrl,
		&i.GitlabUrl,
		&i.Insecure,
		&i.Interval,
		&i.AccessToken,
	)
	return i, err
}

const insertClientRepo = `-- name: InsertClientRepo :exec
INSERT OR IGNORE INTO gitlab_repos
    (name, repo_id, client_id, repo_web_url)
VALUES (?, ?, ?, ?)
`

type InsertClientRepoParams struct {
	Name       string `json:"name"`
	RepoID     int64  `json:"repo_id"`
	ClientID   int64  `json:"client_id"`
	RepoWebUrl string `json:"repo_web_url"`
}

func (q *Queries) InsertClientRepo(ctx context.Context, arg InsertClientRepoParams) error {
	_, err := q.db.ExecContext(ctx, insertClientRepo,
		arg.Name,
		arg.RepoID,
		arg.ClientID,
		arg.RepoWebUrl,
	)
	return err
}

const insertNewClient = `-- name: InsertNewClient :one
INSERT INTO gitlab_clients
(name, created_by, access_token, webhook_url, gitlab_url, insecure)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id
`

type InsertNewClientParams struct {
	Name        string `json:"name"`
	CreatedBy   int64  `json:"created_by"`
	AccessToken string `json:"access_token"`
	WebhookUrl  string `json:"webhook_url"`
	GitlabUrl   string `json:"gitlab_url"`
	Insecure    string `json:"insecure"`
}

func (q *Queries) InsertNewClient(ctx context.Context, arg InsertNewClientParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertNewClient,
		arg.Name,
		arg.CreatedBy,
		arg.AccessToken,
		arg.WebhookUrl,
		arg.GitlabUrl,
		arg.Insecure,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateExistingClient = `-- name: UpdateExistingClient :one
UPDATE gitlab_clients
SET name         = ?,
    created_by   = ?,
    access_token = ?,
    gitlab_url   = ?,
    insecure     = ?
WHERE id = ?
RETURNING id
`

type UpdateExistingClientParams struct {
	Name        string `json:"name"`
	CreatedBy   int64  `json:"created_by"`
	AccessToken string `json:"access_token"`
	GitlabUrl   string `json:"gitlab_url"`
	Insecure    string `json:"insecure"`
	ID          int64  `json:"id"`
}

func (q *Queries) UpdateExistingClient(ctx context.Context, arg UpdateExistingClientParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateExistingClient,
		arg.Name,
		arg.CreatedBy,
		arg.AccessToken,
		arg.GitlabUrl,
		arg.Insecure,
		arg.ID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateTrackedRepoStatus = `-- name: UpdateTrackedRepoStatus :exec
UPDATE gitlab_repos
SET tracked = ?
WHERE repo_id = ?
`

type UpdateTrackedRepoStatusParams struct {
	Tracked int64 `json:"tracked"`
	RepoID  int64 `json:"repo_id"`
}

func (q *Queries) UpdateTrackedRepoStatus(ctx context.Context, arg UpdateTrackedRepoStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateTrackedRepoStatus, arg.Tracked, arg.RepoID)
	return err
}

const upsertClientRepo = `-- name: UpsertClientRepo :exec
INSERT INTO gitlab_repos
    (name, repo_id, client_id, repo_web_url)
VALUES (?, ?, ?, ?)
ON CONFLICT(repo_id) DO UPDATE SET
    name = excluded.name,
    client_id = excluded.client_id,
    repo_web_url = excluded.repo_web_url
`

type UpsertClientRepoParams struct {
	Name       string `json:"name"`
	RepoID     int64  `json:"repo_id"`
	ClientID   int64  `json:"client_id"`
	RepoWebUrl string `json:"repo_web_url"`
}

// INSERT OR REPLACE INTO gitlab_repos
//
//	(name, repo_id, client_id, repo_web_url)
//
// VALUES (?, ?, ?, ?);
func (q *Queries) UpsertClientRepo(ctx context.Context, arg UpsertClientRepoParams) error {
	_, err := q.db.ExecContext(ctx, upsertClientRepo,
		arg.Name,
		arg.RepoID,
		arg.ClientID,
		arg.RepoWebUrl,
	)
	return err
}
