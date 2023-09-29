-- name: InsertNewClient :one
INSERT INTO gitlab_clients
(name, created_by, access_token,  gitlab_url, insecure)
VALUES (?, ?, ?, ?,  ?)
RETURNING id;

-- name: UpdateExistingClient :one
UPDATE gitlab_clients
SET name         = ?,
    created_by   = ?,
    access_token = ?,
    gitlab_url   = ?,
    insecure     = ?
WHERE id = ?
RETURNING id;

-- name: GetClientById :one
SELECT *
FROM gitlab_clients
WHERE id = ?;

-- name: GetAllClients :many
SELECT gc.*, COUNT(gr.repo_id) AS repo_count
FROM gitlab_clients gc
         LEFT JOIN gitlab_repos gr ON gc.id = gr.client_id
GROUP BY gc.id;

-- name: GetAllClientRepos :many
SELECT gr.*
FROM gitlab_repos gr
         INNER JOIN gitlab_clients gc ON gr.client_id = gc.id
WHERE gc.id = ?;

-- name: InsertClientRepo :exec
INSERT OR IGNORE INTO gitlab_repos
    (name, repo_id, client_id, repo_web_url)
VALUES (?, ?, ?, ?);

-- name: UpsertClientRepo :exec
INSERT INTO gitlab_repos
    (name, repo_id, client_id, repo_web_url)
VALUES (?, ?, ?, ?)
ON CONFLICT(repo_id) DO UPDATE SET name         = excluded.name,
                                   client_id    = excluded.client_id,
                                   repo_web_url = excluded.repo_web_url;

-- name: UpdateTrackedRepoStatus :exec
UPDATE gitlab_repos
SET tracked = ?
WHERE repo_id = ?;

-- name: DeleteClientByID :exec
DELETE FROM gitlab_clients
WHERE id = ?;
