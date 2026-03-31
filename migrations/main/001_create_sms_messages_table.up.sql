CREATE TABLE IF NOT EXISTS sms_messages (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES messages(id),
    phone_number VARCHAR(20) NOT NULL,
    text TEXT NOT NULL,
    provider_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_sms_message_id ON sms_messages(message_id);