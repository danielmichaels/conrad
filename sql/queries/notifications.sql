-- name: InsertNotification :one
INSERT INTO notifications (enabled, name, client_id, ignore_approved,
                           ignore_drafts, remind_authors, ignore_approved,
                           min_age, min_staleness, ignore_terms, ignore_labels,
                           require_labels,
                           monday,
                           tuesday,
                           wednesday,
                           thursday,
                           friday,
                           saturday,
                           sunday)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;
-- name: InsetNotificationTimes :one
INSERT INTO notification_times (notification_id, scheduled_time)
VALUES (?, ?)
RETURNING id;
-- name: InsertNotificationMattermost :one
INSERT INTO notifications_mattermost (id, mattermost_channel, webhook_url)
VALUES (?, ?, ?)
RETURNING id;
