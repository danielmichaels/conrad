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
	GitlabUrl   string `json:"gitlab_url"`
	Insecure    string `json:"insecure"`
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

type NotificationTimes struct {
	ID             int64  `json:"id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	NotificationID int64  `json:"notification_id"`
	ScheduledTime  string `json:"scheduled_time"`
	Timezone       string `json:"timezone"`
}

type Notifications struct {
	ID             int64          `json:"id"`
	Enabled        int64          `json:"enabled"`
	Name           string         `json:"name"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	ClientID       int64          `json:"client_id"`
	IgnoreDrafts   int64          `json:"ignore_drafts"`
	RemindAuthors  int64          `json:"remind_authors"`
	IgnoreApproved int64          `json:"ignore_approved"`
	MinAge         int64          `json:"min_age"`
	MinStaleness   int64          `json:"min_staleness"`
	IgnoreTerms    sql.NullString `json:"ignore_terms"`
	IgnoreLabels   sql.NullString `json:"ignore_labels"`
	RequireLabels  int64          `json:"require_labels"`
	Days           string         `json:"days"`
}

type NotificationsEmail struct {
	ID           int64  `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	EmailAddress string `json:"email_address"`
}

type NotificationsMattermost struct {
	ID                int64  `json:"id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	NotificationID    int64  `json:"notification_id"`
	MattermostChannel string `json:"mattermost_channel"`
	WebhookUrl        string `json:"webhook_url"`
}

type NotificationsSlack struct {
	ID           int64  `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	SlackChannel string `json:"slack_channel"`
}

type Users struct {
	ID             int64          `json:"id"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"hashed_password"`
	Name           sql.NullString `json:"name"`
}
