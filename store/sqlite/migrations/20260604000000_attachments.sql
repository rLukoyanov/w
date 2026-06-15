-- +goose Up
CREATE TABLE attachments (
    id           TEXT PRIMARY KEY,
    channel_id   TEXT NOT NULL REFERENCES channels(id),
    author_id    TEXT NOT NULL REFERENCES users(id),
    filename     TEXT NOT NULL,
    size         INTEGER NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'application/octet-stream',
    storage_path TEXT NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attachments_channel_id ON attachments(channel_id);

-- +goose Down
DROP TABLE attachments;
