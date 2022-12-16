package fetcher

import (
	"context"
	"time"

	"github.com/itiky/bb-telegram-notifs/pkg/logging"
)

// runEventsGC runs the events GC worker.
func (f *Fetcher) runEventsGC(ctx context.Context, thesholdDur time.Duration) {
	_, logger := logging.GetCtxLogger(ctx)

	eventIDs, err := f.storage.ListOutdatedEventIDs(ctx, thesholdDur)
	if err != nil {
		logger.Error().Err(err).Msg("List outdated events")
		return
	}
	if len(eventIDs) == 0 {
		logger.Info().Msg("No outdated events found")
		return
	}

	if err := f.storage.DeleteEventByIDs(ctx, eventIDs); err != nil {
		logger.Error().Err(err).Msg("Delete events")
		return
	}

	logger.Info().
		Int("count", len(eventIDs)).
		Msg("Events GC")
}

// runEventsRetry retries the Telegram send for failed events (not acked).
func (f *Fetcher) runEventsRetry(ctx context.Context) {
	_, logger := logging.GetCtxLogger(ctx)

	events, err := f.storage.ListNonAckedEvents(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list non-acked events")
		return
	}
	if len(events) == 0 {
		logger.Info().Msg("No non-acked events found")
		return
	}

	eventsSent := 0
	for _, e := range events {
		ctx, _ := logging.SetCtxLoggerStrFields(ctx,
			"event_hash", e.Hash,
		)

		if f.sendEvent(ctx, e) {
			eventsSent++
		}
	}

	logger.Info().
		Int("count_successful", eventsSent).
		Int("count_failed", len(events)-eventsSent).
		Msg("Events send retry")
}
