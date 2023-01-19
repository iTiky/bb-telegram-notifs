package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/itiky/bb-telegram-notifs/model"
)

// CreateEventSafe creates a new model.Event.
// Idempotent call: returns non-nil object if the event was created, nil otherwise.
func (st *Psql) CreateEventSafe(ctx context.Context, event model.Event) (*model.Event, error) {
	_, err := st.db.NewInsert().
		Model(&event).
		Returning("*").
		Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			if pgErr.IntegrityViolation() {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("creating event: %w", err)
	}

	return &event, nil
}

// SetEventSendAck sets event send acknowledgement.
func (st *Psql) SetEventSendAck(ctx context.Context, eventID int64) error {
	_, err := st.db.NewUpdate().
		Model(&model.Event{}).
		Set("send_ack = true, send_at = ?", time.Now().UTC()).
		Where("id = ?", eventID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("updating event send ack: %w", err)
	}

	return nil
}

// DeleteEventByIDs deletes events by IDs in batch.
func (st *Psql) DeleteEventByIDs(ctx context.Context, IDs []int64) error {
	_, err := st.db.NewDelete().
		Model(&model.Event{}).
		Where("id IN (?)", bun.In(IDs)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting events: %w", err)
	}

	return nil
}

// ListOutdatedEventIDs returns IDs of outdated events (by SendAt AND acked).
func (st *Psql) ListOutdatedEventIDs(ctx context.Context, thresholdDur time.Duration) ([]int64, error) {
	var events []model.Event
	err := st.db.NewSelect().
		Model(&events).
		Where("send_at < ? AND send_ack = true", time.Now().Add(-thresholdDur)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("selecting events: %w", err)
	}

	eventIDs := make([]int64, len(events))
	for i, event := range events {
		eventIDs[i] = event.ID
	}

	return eventIDs, nil
}

// ListNonAckedEvents returns non-acked events.
func (st *Psql) ListNonAckedEvents(ctx context.Context) ([]model.Event, error) {
	var events []model.Event
	err := st.db.NewSelect().
		Model(&events).
		Where("send_ack = false").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("selecting non-acked events: %w", err)
	}

	return events, nil
}
