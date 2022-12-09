package fetcher

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/itiky/bb-telegram-notifs/pkg/logging"
	"github.com/itiky/bb-telegram-notifs/provider/bitbucket"
	"github.com/itiky/bb-telegram-notifs/provider/storage"
	"github.com/itiky/bb-telegram-notifs/service/telegram"
)

// Fetcher is a service for fetching BitBucket PR activity and generating events for Telegram send.
type Fetcher struct {
	storage   *storage.Psql
	bbClient  *bitbucket.Client
	tgService *telegram.Service
}

// NewFetcher creates a new Fetcher.
// Constructor syncs stored repos with BitBucket.
func NewFetcher(ctx context.Context, bbClient *bitbucket.Client, tgSvc *telegram.Service, st *storage.Psql) (*Fetcher, error) {
	if st == nil {
		return nil, fmt.Errorf("storage is nil")
	}
	if bbClient == nil {
		return nil, fmt.Errorf("bitBucket client is nil")
	}
	if tgSvc == nil {
		return nil, fmt.Errorf("telegram service is nil")
	}

	svc := &Fetcher{
		storage:   st,
		bbClient:  bbClient,
		tgService: tgSvc,
	}
	if err := svc.SyncRepos(ctx); err != nil {
		return nil, fmt.Errorf("syncing repos: %w", err)
	}

	return svc, nil
}

// SyncRepos syncs stored repos with BitBucket.
func (f *Fetcher) SyncRepos(ctx context.Context) error {
	ctx, logger := f.Logger(ctx, "sync_repos")

	bbRepos, err := f.bbClient.ListRepos(ctx)
	if err != nil {
		return fmt.Errorf("fetching BB repos: %w", err)
	}

	for _, bbRepo := range bbRepos {
		bbProject, bbRepo := bbRepo.Project.Key, bbRepo.Slug

		repo, err := f.storage.GetRepo(ctx, bbProject, bbRepo)
		if err != nil {
			return fmt.Errorf("getting repo: %w", err)
		}
		if repo != nil {
			logger.Debug().
				Str("project", bbProject).
				Str("repo", bbRepo).
				Msg("Repo already exists")
			continue
		}

		if err := f.storage.CreateRepo(ctx, bbProject, bbRepo); err != nil {
			return fmt.Errorf("creating repo (project: %s, repo: %s): %w", bbProject, bbRepo, err)
		}
		logger.Info().
			Str("project", bbProject).
			Str("repo", bbRepo).
			Msg("Repo created")
	}

	return nil
}

// Logger returns a logger with a service and operation name.
func (f *Fetcher) Logger(ctx context.Context, op string) (context.Context, zerolog.Logger) {
	return logging.SetCtxLoggerStrFields(ctx,
		logging.KeyService, "bb_fetcher",
		logging.KeyOperation, op,
	)
}
