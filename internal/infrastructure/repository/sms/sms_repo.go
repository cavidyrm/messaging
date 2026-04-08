package sms

import (
	"context"
	"database/sql"
	"fmt"
	"messaging/internal/domain/sms"

	"github.com/google/uuid"
)

type SMSRepository struct {
	db *sql.DB
}

func NewSMSRepository(db *sql.DB) *SMSRepository {
	return &SMSRepository{db: db}
}

func (r *SMSRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE sms_messages
		SET status = $1, updated_at = NOW(), sent_at = CASE WHEN $1 = 'sent' THEN NOW() ELSE sent_at END
		WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *SMSRepository) FindByID(ctx context.Context, id uuid.UUID) (*sms.Message, error) {
	var msg sms.Message
	err := r.db.QueryRowContext(ctx, `SELECT id, phone_number, text, status, created_at, updated_at FROM sms_messages WHERE id = $1`, id).
		Scan(&msg.ID, msg.PhoneNumber, msg.Text, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf(`message with id %s not found`, id)
	}
	if err != nil {
		return nil, fmt.Errorf("query sms by id: %w", err)
	}
	return &msg, nil
}

func (r *SMSRepository) Save(ctx context.Context, m *sms.Message) error {
	query := `
		INSERT INTO sms_messages (id, phone_number, text, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, m.ID, m.PhoneNumber, m.Text, m.Status, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert sms failed: %w", err)
	}

	return nil
}
