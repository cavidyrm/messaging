CREATE TABLE IF NOT EXISTS email_notifications (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL,
    to_address  TEXT        NOT NULL,
    subject     TEXT        NOT NULL DEFAULT '',
    body        TEXT        NOT NULL,
    attachments TEXT[]      NOT NULL DEFAULT '{}',
    status      TEXT        NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_user_id ON email_notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_email_status  ON email_notifications (status);
