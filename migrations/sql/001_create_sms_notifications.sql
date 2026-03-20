CREATE TABLE IF NOT EXISTS sms_notifications (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL,
    phone_number TEXT       NOT NULL,
    body        TEXT        NOT NULL,
    status      TEXT        NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sms_user_id  ON sms_notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_sms_status   ON sms_notifications (status);
