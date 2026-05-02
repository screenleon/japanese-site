-- 0015_read_log.sql
-- Per-content read tracking for JS-017. Single-user (no user_id);
-- the DB itself is the user identity for the local-mode deployment.
-- Cloud deployments use NullProgressStore and never touch this table.
-- See BACKLOG.md#JS-017.

CREATE TABLE read_log (
    content_type    TEXT NOT NULL,
    slug            TEXT NOT NULL,
    first_read_at   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_read_at    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_count      INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (content_type, slug)
);

CREATE INDEX idx_read_log_type ON read_log(content_type);
