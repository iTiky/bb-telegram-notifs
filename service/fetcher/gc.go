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

	if err := f.storage.DeleteEventByIDs(ctx, eventIDs); err != nil {
		logger.Error().Err(err).Msg("Delete events")
		return
	}

	logger.Info().Int("count", len(eventIDs)).Msg("Events GC")
}
