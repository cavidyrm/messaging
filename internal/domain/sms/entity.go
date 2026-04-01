package sms

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"
)

type Message struct {
	ID          uuid.UUID
	PhoneNumber string
	Text        string
	Status      string
	SentAt      time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSMS(phoneNumber, text string) *Message {
	return &Message{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		Text:        text,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type Repository interface {
	Save(ctx context.Context, sms *Message) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, sentAt *time.Time) error
	FindByID(ctx context.Context, id uuid.UUID) (*Message, error)
}

type EventRepository interface {
	SaveEvent(ctx context.Context, sms *Message) error
}

type Sender interface {
	Send(ctx context.Context, phoneNumber, text string) error
}
