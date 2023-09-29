// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: notifications.sql

package repository

import (
	"context"
	"database/sql"
)

const deleteNotificationTimesByID = `-- name: DeleteNotificationTimesByID :exec
DELETE FROM notification_times
WHERE notification_id = ?
`

func (q *Queries) DeleteNotificationTimesByID(ctx context.Context, notificationID int64) error {
	_, err := q.db.ExecContext(ctx, deleteNotificationTimesByID, notificationID)
	return err
}

const deleteNotificationsByID = `-- name: DeleteNotificationsByID :exec
DELETE FROM notifications
WHERE id = ?
`

func (q *Queries) DeleteNotificationsByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteNotificationsByID, id)
	return err
}

const deleteNotificationsMattermostByID = `-- name: DeleteNotificationsMattermostByID :exec
DELETE FROM notifications_mattermost
WHERE notification_id = ?
`

func (q *Queries) DeleteNotificationsMattermostByID(ctx context.Context, notificationID int64) error {
	_, err := q.db.ExecContext(ctx, deleteNotificationsMattermostByID, notificationID)
	return err
}

const getAllNotifications = `-- name: GetAllNotifications :many
SELECT n.id, n.enabled, n.name, n.created_at, n.updated_at, n.client_id, n.ignore_drafts, n.remind_authors, n.ignore_approved, n.min_age, n.min_staleness, n.ignore_terms, n.ignore_labels, n.require_labels, n.days,
       nt.id AS notification_time_id,
       nt.scheduled_time,
       mm.id as mattermost_id,
       mm.mattermost_channel,
       mm.webhook_url
FROM notifications n
         LEFT JOIN
     notification_times nt
     ON n.id = nt.notification_id
         LEFT JOIN
     notifications_mattermost mm
     ON n.id = mm.notification_id
`

type GetAllNotificationsRow struct {
	ID                 int64          `json:"id"`
	Enabled            int64          `json:"enabled"`
	Name               string         `json:"name"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
	ClientID           int64          `json:"client_id"`
	IgnoreDrafts       int64          `json:"ignore_drafts"`
	RemindAuthors      int64          `json:"remind_authors"`
	IgnoreApproved     int64          `json:"ignore_approved"`
	MinAge             int64          `json:"min_age"`
	MinStaleness       int64          `json:"min_staleness"`
	IgnoreTerms        sql.NullString `json:"ignore_terms"`
	IgnoreLabels       sql.NullString `json:"ignore_labels"`
	RequireLabels      int64          `json:"require_labels"`
	Days               string         `json:"days"`
	NotificationTimeID sql.NullInt64  `json:"notification_time_id"`
	ScheduledTime      sql.NullString `json:"scheduled_time"`
	MattermostID       sql.NullInt64  `json:"mattermost_id"`
	MattermostChannel  sql.NullString `json:"mattermost_channel"`
	WebhookUrl         sql.NullString `json:"webhook_url"`
}

func (q *Queries) GetAllNotifications(ctx context.Context) ([]GetAllNotificationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllNotifications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllNotificationsRow{}
	for rows.Next() {
		var i GetAllNotificationsRow
		if err := rows.Scan(
			&i.ID,
			&i.Enabled,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ClientID,
			&i.IgnoreDrafts,
			&i.RemindAuthors,
			&i.IgnoreApproved,
			&i.MinAge,
			&i.MinStaleness,
			&i.IgnoreTerms,
			&i.IgnoreLabels,
			&i.RequireLabels,
			&i.Days,
			&i.NotificationTimeID,
			&i.ScheduledTime,
			&i.MattermostID,
			&i.MattermostChannel,
			&i.WebhookUrl,
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

const getNotificationByID = `-- name: GetNotificationByID :one
SELECT n.id, n.enabled, n.name, n.created_at, n.updated_at, n.client_id, n.ignore_drafts, n.remind_authors, n.ignore_approved, n.min_age, n.min_staleness, n.ignore_terms, n.ignore_labels, n.require_labels, n.days,
       nt.id AS notification_time_id,
       nt.scheduled_time,
       mm.id as mattermost_id,
       mm.mattermost_channel,
       mm.webhook_url
FROM notifications n
         LEFT JOIN
     notification_times nt
     ON n.id = nt.notification_id
         LEFT JOIN
     notifications_mattermost mm
     ON n.id = mm.notification_id
WHERE n.id = ?
`

type GetNotificationByIDRow struct {
	ID                 int64          `json:"id"`
	Enabled            int64          `json:"enabled"`
	Name               string         `json:"name"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
	ClientID           int64          `json:"client_id"`
	IgnoreDrafts       int64          `json:"ignore_drafts"`
	RemindAuthors      int64          `json:"remind_authors"`
	IgnoreApproved     int64          `json:"ignore_approved"`
	MinAge             int64          `json:"min_age"`
	MinStaleness       int64          `json:"min_staleness"`
	IgnoreTerms        sql.NullString `json:"ignore_terms"`
	IgnoreLabels       sql.NullString `json:"ignore_labels"`
	RequireLabels      int64          `json:"require_labels"`
	Days               string         `json:"days"`
	NotificationTimeID sql.NullInt64  `json:"notification_time_id"`
	ScheduledTime      sql.NullString `json:"scheduled_time"`
	MattermostID       sql.NullInt64  `json:"mattermost_id"`
	MattermostChannel  sql.NullString `json:"mattermost_channel"`
	WebhookUrl         sql.NullString `json:"webhook_url"`
}

func (q *Queries) GetNotificationByID(ctx context.Context, id int64) (GetNotificationByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getNotificationByID, id)
	var i GetNotificationByIDRow
	err := row.Scan(
		&i.ID,
		&i.Enabled,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ClientID,
		&i.IgnoreDrafts,
		&i.RemindAuthors,
		&i.IgnoreApproved,
		&i.MinAge,
		&i.MinStaleness,
		&i.IgnoreTerms,
		&i.IgnoreLabels,
		&i.RequireLabels,
		&i.Days,
		&i.NotificationTimeID,
		&i.ScheduledTime,
		&i.MattermostID,
		&i.MattermostChannel,
		&i.WebhookUrl,
	)
	return i, err
}

const insertNotification = `-- name: InsertNotification :one
INSERT INTO notifications
(enabled, name, client_id, ignore_approved, ignore_drafts, remind_authors,
 min_age, min_staleness, ignore_terms, ignore_labels, require_labels, days)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id
`

type InsertNotificationParams struct {
	Enabled        int64          `json:"enabled"`
	Name           string         `json:"name"`
	ClientID       int64          `json:"client_id"`
	IgnoreApproved int64          `json:"ignore_approved"`
	IgnoreDrafts   int64          `json:"ignore_drafts"`
	RemindAuthors  int64          `json:"remind_authors"`
	MinAge         int64          `json:"min_age"`
	MinStaleness   int64          `json:"min_staleness"`
	IgnoreTerms    sql.NullString `json:"ignore_terms"`
	IgnoreLabels   sql.NullString `json:"ignore_labels"`
	RequireLabels  int64          `json:"require_labels"`
	Days           string         `json:"days"`
}

func (q *Queries) InsertNotification(ctx context.Context, arg InsertNotificationParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertNotification,
		arg.Enabled,
		arg.Name,
		arg.ClientID,
		arg.IgnoreApproved,
		arg.IgnoreDrafts,
		arg.RemindAuthors,
		arg.MinAge,
		arg.MinStaleness,
		arg.IgnoreTerms,
		arg.IgnoreLabels,
		arg.RequireLabels,
		arg.Days,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateEnabledNotificationStatus = `-- name: UpdateEnabledNotificationStatus :exec
UPDATE notifications
SET enabled = ?
WHERE notifications.id = ?
`

type UpdateEnabledNotificationStatusParams struct {
	Enabled int64 `json:"enabled"`
	ID      int64 `json:"id"`
}

func (q *Queries) UpdateEnabledNotificationStatus(ctx context.Context, arg UpdateEnabledNotificationStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateEnabledNotificationStatus, arg.Enabled, arg.ID)
	return err
}

const updateNotificationMattermost = `-- name: UpdateNotificationMattermost :one
UPDATE notifications_mattermost
SET
    mattermost_channel  = ?,
    webhook_url         = ?
WHERE notification_id = ?
RETURNING id
`

type UpdateNotificationMattermostParams struct {
	MattermostChannel string `json:"mattermost_channel"`
	WebhookUrl        string `json:"webhook_url"`
	NotificationID    int64  `json:"notification_id"`
}

func (q *Queries) UpdateNotificationMattermost(ctx context.Context, arg UpdateNotificationMattermostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateNotificationMattermost, arg.MattermostChannel, arg.WebhookUrl, arg.NotificationID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateNotificationTimes = `-- name: UpdateNotificationTimes :one
UPDATE notification_times
SET
    scheduled_time   = ?,
    timezone         = ?
WHERE notification_id = ?
RETURNING id
`

type UpdateNotificationTimesParams struct {
	ScheduledTime  string `json:"scheduled_time"`
	Timezone       string `json:"timezone"`
	NotificationID int64  `json:"notification_id"`
}

func (q *Queries) UpdateNotificationTimes(ctx context.Context, arg UpdateNotificationTimesParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateNotificationTimes, arg.ScheduledTime, arg.Timezone, arg.NotificationID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const upsertNotification = `-- name: UpsertNotification :one
INSERT INTO notifications
(id,
 enabled,
 name,
 client_id,
 ignore_approved,
 ignore_drafts,
 remind_authors,
 min_age,
 min_staleness,
 ignore_terms,
 ignore_labels,
 require_labels,
 days)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET enabled         = excluded.enabled,
                              name            = excluded.name,
                              ignore_approved = excluded.ignore_approved,
                              ignore_drafts   = excluded.ignore_drafts,
                              remind_authors  = excluded.remind_authors,
                              min_age         = excluded.min_age,
                              min_staleness   = excluded.min_staleness,
                              ignore_terms    = excluded.ignore_terms,
                              ignore_labels   = excluded.ignore_labels,
                              require_labels  = excluded.require_labels,
                              days            = excluded.days
RETURNING id
`

type UpsertNotificationParams struct {
	ID             int64          `json:"id"`
	Enabled        int64          `json:"enabled"`
	Name           string         `json:"name"`
	ClientID       int64          `json:"client_id"`
	IgnoreApproved int64          `json:"ignore_approved"`
	IgnoreDrafts   int64          `json:"ignore_drafts"`
	RemindAuthors  int64          `json:"remind_authors"`
	MinAge         int64          `json:"min_age"`
	MinStaleness   int64          `json:"min_staleness"`
	IgnoreTerms    sql.NullString `json:"ignore_terms"`
	IgnoreLabels   sql.NullString `json:"ignore_labels"`
	RequireLabels  int64          `json:"require_labels"`
	Days           string         `json:"days"`
}

func (q *Queries) UpsertNotification(ctx context.Context, arg UpsertNotificationParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, upsertNotification,
		arg.ID,
		arg.Enabled,
		arg.Name,
		arg.ClientID,
		arg.IgnoreApproved,
		arg.IgnoreDrafts,
		arg.RemindAuthors,
		arg.MinAge,
		arg.MinStaleness,
		arg.IgnoreTerms,
		arg.IgnoreLabels,
		arg.RequireLabels,
		arg.Days,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const upsertNotificationMattermost = `-- name: UpsertNotificationMattermost :one
INSERT INTO notifications_mattermost (mattermost_channel, webhook_url,
                                      notification_id)
VALUES (?, ?, ?)
ON CONFLICT(notification_id) DO UPDATE SET notification_id    = excluded.notification_id,
                                           mattermost_channel = excluded.mattermost_channel,
                                           webhook_url        = excluded.webhook_url
RETURNING id
`

type UpsertNotificationMattermostParams struct {
	MattermostChannel string `json:"mattermost_channel"`
	WebhookUrl        string `json:"webhook_url"`
	NotificationID    int64  `json:"notification_id"`
}

func (q *Queries) UpsertNotificationMattermost(ctx context.Context, arg UpsertNotificationMattermostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, upsertNotificationMattermost, arg.MattermostChannel, arg.WebhookUrl, arg.NotificationID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const upsertNotificationTimes = `-- name: UpsertNotificationTimes :one
INSERT INTO notification_times (notification_id, scheduled_time, timezone)
VALUES (?, ?, ?)
ON CONFLICT(notification_id) DO UPDATE SET notification_id = excluded.notification_id,
                                           scheduled_time  = excluded.scheduled_time,
                                           timezone        = excluded.timezone
RETURNING id
`

type UpsertNotificationTimesParams struct {
	NotificationID int64  `json:"notification_id"`
	ScheduledTime  string `json:"scheduled_time"`
	Timezone       string `json:"timezone"`
}

func (q *Queries) UpsertNotificationTimes(ctx context.Context, arg UpsertNotificationTimesParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, upsertNotificationTimes, arg.NotificationID, arg.ScheduledTime, arg.Timezone)
	var id int64
	err := row.Scan(&id)
	return id, err
}
