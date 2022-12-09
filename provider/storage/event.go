package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/itiky/bb-telegram-notifs/model"
)

// CreateEventSafe creates a new model.Event.
// Idempotent call: returns true if the event was created, false if the event already exists.
func (st *Psql) CreateEventSafe(ctx context.Context, event model.Event) (bool, error) {
	_, err := st.db.NewInsert().
		Model(&event).
		Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			if pgErr.IntegrityViolation() {
				return false, nil
			}
		}
		return false, fmt.Errorf("creating event: %w", err)
	}

	return true, nil
}

// DeleteEventByIDs deletes events by IDs in batch.
func (st *Psql) DeleteEventByIDs(ctx context.Context, IDs []int64) error {
	_, err := st.db.NewDelete().
		Model(&model.Event{}).
		Where("id IN (?)", IDs).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting events: %w", err)
	}

	return nil
}

// ListOutdatedEventIDs returns IDs of outdated events (by CreatedAt).
func (st *Psql) ListOutdatedEventIDs(ctx context.Context, thresholdDur time.Duration) ([]int64, error) {
	var events []model.Event
	err := st.db.NewSelect().
		Model(&events).
		Where("created_at < ?", time.Now().Add(-thresholdDur)).
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
