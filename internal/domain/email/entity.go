package email

import (
	"context"
	"messaging/internal/domain/event"

	"github.com/google/uuid"
)

type Email struct {
	Address string
	Subject string
	Body    string
}

type Repository interface {
	Save(ctx context.Context, msg *Email) error
	UpdateStatus(ctx context.Context, is uuid.UUID, status string) error
}

type EventRecorder interface {
	RecordStateChange(ctx context.Context, event event.Event, snapshot event.Snapshot) error
}

type Sender interface {
	Send(ctx context.Context, email *Email) error
}
