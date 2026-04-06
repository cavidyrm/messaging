package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EventID       uuid.UUID       `json:"event_id"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"` // "SMS", "Email", "Push"
	EventType     string          `json:"event_type"`     // "SMSSendRequested"
	Version       int             `json:"version"`
	Data          json.RawMessage `json:"data"`
	Metadata      json.RawMessage `json:"metadata"`
	Timestamp     time.Time       `json:"timestamp"`
}

type Snapshot struct {
	SnapshotID    uuid.UUID       `json:"snapshot_id"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	Version       int             `json:"version"`
	State         json.RawMessage `json:"state"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type Repository interface {
	SaveEventAndSnapshot(ctx context.Context, event Event, snapshot Snapshot) error
}
