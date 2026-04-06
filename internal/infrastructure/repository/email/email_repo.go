package email

import (
	"context"
	"database/sql"
	"messaging/internal/domain/email"

	"github.com/google/uuid"
)

type EmailRepository struct {
	db *sql.DB
}

func NewEmailRepository(db *sql.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

func (r *EmailRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	//TODO implement me
	return nil
}

func (r *EmailRepository) FindByID(ctx context.Context, id uuid.UUID) (*email.Email, error) {
	//TODO implement me
	return nil, nil
}

func (r *EmailRepository) Save(ctx context.Context, m *email.Email) error {
	//tx, err := r.db.BeginTx(ctx, nil)
	//{
	//	if err != nil {
	//		return fmt.Errorf("transaction begin: %w", err)
	//	}
	//	defer tx.Rollback()
	//}
	//
	//query := `INSERT INTO sms_messages (id, phone_number, text, status, created_at, updated_at)  VALUES ($1, $2, $3, $4, $5, $6)`
	//
	//_, err = tx.ExecContext(ctx, query, m.ID, m.PhoneNumber, m.Text, m.Status, m.CreatedAt, m.UpdatedAt)
	//if err != nil {
	//	return fmt.Errorf("insert sms failed: %w", err)
	//}
	//err = tx.Commit()
	//if err != nil {
	//	return fmt.Errorf("commit transaction failed: %w", err)
	//}
	return nil
}
