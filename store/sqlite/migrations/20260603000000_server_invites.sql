-- +goose Up
CREATE TABLE server_invites (
    id          TEXT PRIMARY KEY,
    server_id   TEXT NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    created_by  TEXT NOT NULL REFERENCES users(id),
    code        TEXT NOT NULL UNIQUE,
    max_uses    INTEGER,
    use_count   INTEGER NOT NULL DEFAULT 0,
    expires_at  DATETIME,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE server_invites;
