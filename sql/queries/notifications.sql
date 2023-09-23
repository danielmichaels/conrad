-- name: GetAllNotifications :many
SELECT n.*,
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
     ON n.id = mm.id;

-- name: GetNotificationByID :one
SELECT n.*,
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
     ON n.id = mm.id
WHERE n.id = ?;

-- name: InsertNotification :one
INSERT INTO notifications
(enabled, name, client_id, ignore_approved, ignore_drafts, remind_authors,
 min_age, min_staleness, ignore_terms, ignore_labels, require_labels, days)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: UpsertNotification :one
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
RETURNING id;

-- name: UpsertNotificationTimes :one
INSERT INTO notification_times (notification_id, scheduled_time, timezone)
VALUES (?, ?, ?)
ON CONFLICT(notification_id) DO UPDATE SET notification_id = excluded.notification_id,
                                           scheduled_time  = excluded.scheduled_time,
                                           timezone        = excluded.timezone
RETURNING id;

-- name: UpsertNotificationMattermost :one
INSERT INTO notifications_mattermost (mattermost_channel, webhook_url,
                                      notification_id)
VALUES (?, ?, ?)
ON CONFLICT(notification_id) DO UPDATE SET notification_id    = excluded.notification_id,
                                           mattermost_channel = excluded.mattermost_channel,
                                           webhook_url        = excluded.webhook_url
RETURNING id;

-- name: UpdateEnabledNotificationStatus :exec
UPDATE notifications
SET enabled = ?
WHERE notifications.id = ?;
