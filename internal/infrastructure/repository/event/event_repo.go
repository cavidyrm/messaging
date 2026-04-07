package event

import (
	"context"
	"database/sql"
	"fmt"

	"messaging/internal/domain/event"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventStore(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// SaveEventAndSnapshot saves both records atomically to avoid dual-write issues
func (es *EventRepository) SaveEvent(ctx context.Context, ev event.Event) error {
	tx, err := es.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insert Event
	insertEventQuery := `
		INSERT INTO microservices.events (event_id, aggregate_id, aggregate_type, event_type, version, data, metadata, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, insertEventQuery,
		ev.EventID, ev.AggregateID, ev.AggregateType, ev.EventType,
		ev.Version, ev.Data, ev.Metadata, ev.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	// 2. Upsert Snapshot (Insert or update if version is higher)
	//upsertSnapshotQuery := `
	//	INSERT INTO microservices.snapshots (aggregate_id, aggregate_type, version, data, metadata, timestamp)
	//	VALUES ($1, $2, $3, $4, $5, $6)
	//	ON CONFLICT (aggregate_id) DO UPDATE
	//	SET version = EXCLUDED.version,
	//	    updated_at = EXCLUDED.updated_at
	//	WHERE snapshots.version < EXCLUDED.version`
	//
	//_, err = tx.ExecContext(ctx, upsertSnapshotQuery,
	//	snap.AggregateID, snap.AggregateType, snap.Version,
	//)
	//if err != nil {
	//	return fmt.Errorf("failed to upsert snapshot: %w", err)
	//}

	return tx.Commit()
}
