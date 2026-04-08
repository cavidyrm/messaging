package email

import (
	"context"
	"messaging/internal/domain/event"
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID        uuid.UUID `json:"id"`
	Address   string    `json:"address"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Save(ctx context.Context, msg *Email) error
	UpdateStatus(ctx context.Context, is uuid.UUID, status string) error
	FindByID(ctx context.Context, id uuid.UUID) (*Email, error)
}

type EventRecorder interface {
	RecordStateChange(ctx context.Context, event event.Event, snapshot event.Snapshot) error
}

type Sender interface {
	Send(ctx context.Context, email *Email) error
}
