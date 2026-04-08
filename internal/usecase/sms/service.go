package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"messaging/internal/domain/event"
	"time"

	"messaging/internal/domain/sms"

	"github.com/google/uuid"
)

type Service struct {
	eventRepo event.Repository
	smsRepo   sms.Repository
	//logger    *logger.Logger
	sender sms.Sender
}

func NewSMSService(eventRepo event.Repository, smsRepo sms.Repository, sender sms.Sender) *Service {
	return &Service{
		eventRepo: eventRepo,
		smsRepo:   smsRepo,
		sender:    sender,
	}
}

func (s *Service) ProcessAndSendSMS(ctx context.Context, reqEvent event.Event) error {
	now := time.Now()

	var smsData struct {
		PhoneNumber string `json:"phone_number"`
		Text        string `json:"text"`
	}
	if err := json.Unmarshal(reqEvent.Data, &smsData); err != nil {
		return fmt.Errorf("failed to unpack sms data: %w", err)
	}

	// We use the AggregateID and Version provided by the Kafka event
	aggregateID := reqEvent.AggregateID
	currentVersion := reqEvent.Version

	// Create the SMS entity for our Reporting DB
	smsMsg := &sms.Message{
		ID:          aggregateID,
		PhoneNumber: smsData.PhoneNumber,
		Text:        smsData.Text,
		Status:      "PENDING",
		CreatedAt:   reqEvent.Timestamp,
		UpdatedAt:   now,
	}
	//stateBytes, _ := json.Marshal(smsMsg)
	//
	//// Create initial snapshot from the incoming event
	//snapRequested := event.Snapshot{
	//	SnapshotID:    reqEvent.EventID, // Or generate a new UUID
	//	AggregateID:   aggregateID,
	//	AggregateType: "SMS",
	//	Version:       currentVersion,
	//
	//}

	// Save the incoming event and snapshot to the Event Store
	if err := s.eventRepo.SaveEvent(ctx, reqEvent); err != nil {
		return err
	}
	// Save to Reporting DB
	if err := s.smsRepo.Save(ctx, smsMsg); err != nil {
		return err
	}

	// send sms
	sendErr := s.sender.Send(ctx, smsData.PhoneNumber, smsData.Text)

	// sms sent or fialed
	now = time.Now()
	smsMsg.UpdatedAt = now
	currentVersion++ // Increment version

	var eventType string
	var metadataBytes json.RawMessage

	if sendErr != nil {
		smsMsg.Status = "FAILED"
		eventType = "SMSFailed"
		metadataBytes = json.RawMessage(`{"error": "` + sendErr.Error() + `"}`)
	} else {
		smsMsg.Status = "SENT"
		eventType = "SMSSent"
		metadataBytes = json.RawMessage(`{}`)
	}

	newStateBytes, _ := json.Marshal(smsMsg)

	// Create the Outcome Event
	evOutcome := event.Event{
		EventID:       uuid.New(),
		AggregateID:   aggregateID, // Same aggregate!
		AggregateType: "SMS",
		EventType:     eventType,
		Version:       currentVersion, // Incremented Version
		Data:          newStateBytes,
		Metadata:      metadataBytes,
		Timestamp:     now,
	}

	//// Create the Outcome Snapshot
	//snapOutcome := event.Snapshot{
	//	SnapshotID:    uuid.New(),
	//	AggregateID:   aggregateID,
	//	AggregateType: "SMS",
	//	Version:       currentVersion, // Incremented Version
	//	Data:          newStateBytes,
	//	Metadata:      metadataBytes,
	//	Timestamp:     now,
	//}

	// Save Outcome Event & Snapshot
	if err := s.eventRepo.SaveEvent(ctx, evOutcome); err != nil {
		return err
	}
	// Update Reporting DB
	if err := s.smsRepo.UpdateStatus(ctx, aggregateID, smsMsg.Status); err != nil {
		return err
	}

	return sendErr
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*sms.Message, error) {
	return s.smsRepo.FindByID(ctx, id)
}
