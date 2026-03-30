CREATE TABLE IF NOT EXISTS email_messages (
                                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES messages(id),
    email_address VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    text TEXT NOT NULL,
    provider_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_email_message_id ON email_messages(message_id);