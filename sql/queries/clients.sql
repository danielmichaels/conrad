-- name: InsertNewClient :one
INSERT INTO gitlab_clients
(name, created_by, access_token, webhook_url, gitlab_url)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: GetClientById :one
SELECT *
FROM gitlab_clients
WHERE id = ?;

-- name: GetAllClients :many
SELECT gc.*, COUNT(gr.id) AS repo_count
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

-- name: UpdateTrackedRepoStatus :exec
UPDATE gitlab_repos
SET tracked = ?
WHERE repo_id = ?;
