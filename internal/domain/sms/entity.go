package sms

import (
	"context"
	"messaging/internal/domain/event"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending Status = "PENDING"
	StatusSent    Status = "SENT"
	StatusFailed  Status = "FAILED"
)

type Message struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Text        string    `json:"text"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Repository interface {
	Save(ctx context.Context, msg *Message) error
	UpdateStatus(ctx context.Context, is uuid.UUID, status string) error
}

type EventRecorder interface {
	RecordStateChange(ctx context.Context, event event.Event, snapshot event.Snapshot) error
}

type Sender interface {
	Send(ctx context.Context, phoneNumber, text string) error
}
