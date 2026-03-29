CREATE TABLE IF NOT EXISTS push_notifications (
    id            TEXT        PRIMARY KEY,
    user_id       TEXT        NOT NULL,
    device_token  TEXT        NOT NULL,
    title         TEXT        NOT NULL,
    body          TEXT        NOT NULL,
    template_name TEXT        NOT NULL DEFAULT '',
    data          JSONB       NOT NULL DEFAULT '{}',
    status        TEXT        NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_push_user_id ON push_notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_push_status  ON push_notifications (status);
