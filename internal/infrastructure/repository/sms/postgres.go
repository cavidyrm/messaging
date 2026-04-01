package repository

import (
	"context"
	"database/sql"
	"fmt"
	"messaging/internal/domain/sms"
)

type SMSRepository struct {
	db *sql.DB
}

func NewSMSRepository(db *sql.DB) *SMSRepository {
	return &SMSRepository{db: db}
}

func (r *SMSRepository) Save(ctx context.Context, m sms.Message) error {
	tx, err := r.db.BeginTx(ctx, nil)
	{
		if err != nil {
			return fmt.Errorf("transaction begin: %w", err)
		}
		defer tx.Rollback()
	}

	query := `INSERT INTO sms_messages (id, phone_number, text, status, created_at, updated_at)  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, query, m.ID, m.PhoneNumber, m.Text, m.Status, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert sms failed: %w", err)
	}
	err := tx.Commit()
	return nil
}
