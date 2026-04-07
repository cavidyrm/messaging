package email

import (
	"context"
	"database/sql"
	"fmt"
	"messaging/internal/domain/email"

	"github.com/google/uuid"
)

type EmailRepository struct {
	db *sql.DB
}

func NewEmailRepository(db *sql.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

func (r *EmailRepository) Save(ctx context.Context, m *email.Email) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO email_messages (id, to_address, subject, body, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		m.ID, m.Address, m.Subject, m.Body, m.Status, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *EmailRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE email_messages
		 SET status = $1, updated_at = NOW(),
		     sent_at = CASE WHEN status = 'SENT' THEN NOW() ELSE sent_at END
		 WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update email status failed: %w", err)
	}
	return nil
}
