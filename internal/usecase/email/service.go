package email

import (
	"context"
	"encoding/json"
	"fmt"
	"messaging/internal/domain/email"
	"messaging/internal/domain/event"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	eventRepo event.Repository
	emailRepo email.Repository
	sender    email.Sender
}

func NewEmailService(eventRepo event.Repository, emailRepo email.Repository, sender email.Sender) *Service {
	return &Service{
		eventRepo: eventRepo,
		emailRepo: emailRepo,
		sender:    sender,
	}
}

func (s *Service) ProcessAndSendEmail(ctx context.Context, reqEvent event.Event) error {
	now := time.Now()

	var data struct {
		Address string `json:"to_address"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := json.Unmarshal(reqEvent.Data, &data); err != nil {
		return fmt.Errorf("failed to unpack email data: %w", err)
	}

	aggregateID := reqEvent.AggregateID
	currentVersion := reqEvent.Version

	emailMsg := &email.Email{
		ID:        aggregateID,
		Address:   data.Address,
		Subject:   data.Subject,
		Body:      data.Body,
		Status:    "PENDING",
		CreatedAt: reqEvent.Timestamp,
		UpdatedAt: now,
	}

	if err := s.eventRepo.SaveEvent(ctx, reqEvent); err != nil {
		return err
	}
	if err := s.emailRepo.Save(ctx, emailMsg); err != nil {
		return err
	}

	sendErr := s.sender.Send(ctx, emailMsg)

	now = time.Now()
	emailMsg.UpdatedAt = now
	currentVersion++

	var eventType string
	var metadataBytes json.RawMessage

	if sendErr != nil {
		emailMsg.Status = "FAILED"
		eventType = "EmailFailed"
		metadataBytes = json.RawMessage(`{"error": "` + sendErr.Error() + `"}`)
	} else {
		emailMsg.Status = "SENT"
		eventType = "EmailSent"
		metadataBytes = json.RawMessage(`{}`)
	}

	newStateBytes, _ := json.Marshal(emailMsg)

	evOutcome := event.Event{
		EventID:       uuid.New(),
		AggregateID:   aggregateID,
		AggregateType: "Email",
		EventType:     eventType,
		Version:       currentVersion,
		Data:          newStateBytes,
		Metadata:      metadataBytes,
		Timestamp:     now,
	}

	if err := s.eventRepo.SaveEvent(ctx, evOutcome); err != nil {
		return err
	}
	if err := s.emailRepo.UpdateStatus(ctx, aggregateID, emailMsg.Status); err != nil {
		return err
	}

	return sendErr
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*email.Email, error) {
	return s.emailRepo.FindByID(ctx, id)
}
