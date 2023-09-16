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
SELECT *
FROM gitlab_clients;

-- name: InsertClientRepo :exec
INSERT OR IGNORE INTO gitlab_repos
(name,repo_id, client_id)
VALUES (?,?,?);
