// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package repository

import (
	"database/sql"
)

type GitlabClients struct {
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
}

type GitlabRepos struct {
	RepoID     int64  `json:"repo_id"`
	RepoWebUrl string `json:"repo_web_url"`
	Tracked    int64  `json:"tracked"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	ClientID   int64  `json:"client_id"`
}

type Users struct {
	ID             int64          `json:"id"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"hashed_password"`
	Name           sql.NullString `json:"name"`
}
