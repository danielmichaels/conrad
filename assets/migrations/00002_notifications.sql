-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications
(
    id              INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    -- enabled: 0 is false, 1 is true
    enabled         INTEGER NOT NULL DEFAULT 0,
    name            TEXT    NOT NULL,
    created_at      TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    client_id       INTEGER NOT NULL,
    ignore_drafts   INTEGER NOT NULL DEFAULT 0,
    remind_authors  INTEGER NOT NULL DEFAULT 0,
    ignore_approved INTEGER NOT NULL DEFAULT 0,
    -- min_age: don't review if PR < 2 hours old
    min_age         INTEGER NOT NULL DEFAULT 2,
    -- min_staleness: don't review if last activity more recent that this in hours
    min_staleness   INTEGER NOT NULL DEFAULT 2,
    -- ignore_terms: comma seperated list of terms to ignore if in PR title
    ignore_terms    TEXT,
    -- ignore_labels: PR's with these labels are ignored; comma separated list
    ignore_labels   TEXT,
    -- require_labels: only review if PR has one of the labels from this list
    require_labels  INTEGER NOT NULL DEFAULT 0,
    days TEXT NOT NULL,
    FOREIGN KEY (client_id) REFERENCES gitlab_clients (id) ON DELETE CASCADE
);
CREATE TABLE notification_times
(
    id              INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at      TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notification_id INTEGER NOT NULL UNIQUE,
    scheduled_time  TEXT    NOT NULL,
    timezone        TEXT    NOT NULL,
    FOREIGN KEY (notification_id) REFERENCES notifications (id) ON DELETE CASCADE
);
CREATE TABLE notifications_mattermost
(
    id                 INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at         TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notification_id    INTEGER NOT NULL UNIQUE,
    mattermost_channel TEXT    NOT NULL,
    webhook_url        TEXT    NOT NULL,
    -- Other Mattermost-specific columns
    FOREIGN KEY (id) REFERENCES notifications (id) ON DELETE CASCADE
);
CREATE TABLE notifications_slack
(
    id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    slack_channel TEXT    NOT NULL,
    -- Other Slack-specific columns
    FOREIGN KEY (id) REFERENCES notifications (id) ON DELETE CASCADE
);
CREATE TABLE notifications_email
(
    id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email_address TEXT    NOT NULL,
    -- Other Email-specific columns
    FOREIGN KEY (id) REFERENCES notifications (id) ON DELETE CASCADE
);
CREATE TRIGGER updated_at_notifications
    AFTER UPDATE
    ON notifications
    FOR EACH ROW
BEGIN
    UPDATE notifications
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TRIGGER updated_at_slack
    AFTER UPDATE
    ON notifications_slack
    FOR EACH ROW
BEGIN
    UPDATE notifications_slack
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TRIGGER updated_at_mattermost
    AFTER UPDATE
    ON notifications_mattermost
    FOR EACH ROW
BEGIN
    UPDATE notifications_mattermost
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TRIGGER updated_at_email
    AFTER UPDATE
    ON notifications_email
    FOR EACH ROW
BEGIN
    UPDATE notifications_email
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TRIGGER updated_at_notification_times
    AFTER UPDATE
    ON notification_times
    FOR EACH ROW
BEGIN
    UPDATE notification_times
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS updated_at_notifications;
DROP TRIGGER IF EXISTS updated_at_notification_times;
DROP TRIGGER IF EXISTS updated_at_slack;
DROP TRIGGER IF EXISTS updated_at_email;
DROP TRIGGER IF EXISTS updated_at_mattermost;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS notification_times;
DROP TABLE IF EXISTS notifications_email;
DROP TABLE IF EXISTS notifications_mattermost;
DROP TABLE IF EXISTS notifications_slack;
-- +goose StatementEnd
