package fetcher

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/itiky/bb-telegram-notifs/model"
	"github.com/itiky/bb-telegram-notifs/pkg"
	"github.com/itiky/bb-telegram-notifs/pkg/config"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
	bbModel "github.com/itiky/bb-telegram-notifs/provider/bitbucket/model"
)

// userSet is a container for users.
type userSet map[string]model.User // key: bbEmail

// StartWorkers starts the fetcher workers: activity poller and events GC.
func (f *Fetcher) StartWorkers(ctx context.Context) error {
	ctx, logger := logging.SetCtxLoggerStrFields(ctx,
		logging.KeyService, "bb_fetcher",
	)

	fetchPeriod := viper.GetDuration(config.FetchPeriod)
	if fetchPeriod <= 0 {
		return fmt.Errorf("invalid fetch period: %v", fetchPeriod)
	}

	retryPeriod := viper.GetDuration(config.RetryPeriod)
	if retryPeriod <= 0 {
		return fmt.Errorf("invalid retry period: %v", retryPeriod)
	}

	gcPeriod := viper.GetDuration(config.EventsGCPeriod)
	if gcPeriod <= 0 {
		return fmt.Errorf("invalid events GC period: %v", gcPeriod)
	}

	gcThresholdDur := viper.GetDuration(config.EventGCThreshold)
	if gcThresholdDur <= 0 {
		return fmt.Errorf("invalid event GC threshold: %v", gcThresholdDur)
	}

	logger.Info().
		Dur("fetch_period", fetchPeriod).
		Dur("retry_period", retryPeriod).
		Dur("gc_period", gcPeriod).
		Dur("gc_threshold", gcThresholdDur).
		Msg("Workers started")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(fetchPeriod):
				ctx = pkg.ContextWithCorrelationID(ctx)
				ctx, logger := f.Logger(ctx, "fetch")

				repos, err := f.storage.ListRepos(ctx)
				if err != nil {
					logger.Error().Err(err).Msg("List repos")
					break
				}

				var events []model.Event
				for _, repo := range repos {
					events = append(events, f.fetchRepo(ctx, repo)...)
				}
				f.handleEvents(ctx, events)
			case <-time.Tick(retryPeriod):
				ctx = pkg.ContextWithCorrelationID(ctx)
				ctx, _ := f.Logger(ctx, "events_retry")

				f.runEventsRetry(ctx)
			case <-time.Tick(gcPeriod):
				ctx = pkg.ContextWithCorrelationID(ctx)
				ctx, _ := f.Logger(ctx, "events_gc")

				f.runEventsGC(ctx, gcThresholdDur)
			}
		}
	}()

	return nil
}

// fetchRepo polls all the open PRs of the given repo and returns the events.
func (f *Fetcher) fetchRepo(ctx context.Context, repo model.Repo) (events []model.Event) {
	ctx, logger := logging.SetCtxLoggerStrFields(ctx,
		"project", repo.Project,
		"repo", repo.Name,
	)

	subAllUsers, subROnlyUsers, err := f.getPRRelatedUsers(ctx, repo.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Get open PR related users")
		return
	}
	if len(subAllUsers) == 0 && len(subROnlyUsers) == 0 {
		logger.Debug().Msg("No subscribers (skip)")
		return
	}

	openPRs, err := f.bbClient.ListPRsForRepo(ctx, repo.Name, bbModel.PRStateOpen)
	if err != nil {
		logger.Error().Err(err).Msg("List PRs for repo")
		return nil
	}

	for _, pr := range openPRs {
		ctx, logger := logging.SetCtxLoggerStrFields(ctx,
			"pr_id", strconv.FormatInt(pr.ID, 10),
		)

		prEvents := f.fetchOpenPR(ctx, repo, pr, subAllUsers, subROnlyUsers)
		logger.Info().Int("events", len(prEvents)).Msg("Open PR events fetched")

		events = append(events, prEvents...)
	}

	return
}

// fetchOpenPR polls the given open PR activity and returns the events.
func (f *Fetcher) fetchOpenPR(ctx context.Context, repo model.Repo, pr bbModel.PR, subAllUsers, subROnlyUsers userSet) (events []model.Event) {
	ctx, logger := logging.GetCtxLogger(ctx)

	// Merge subscriber sets removing non-reviewers
	subUsers := make(userSet)
	for _, user := range subAllUsers {
		subUsers[user.BbEmail] = user
	}
	for _, bbUser := range pr.Reviewers {
		user, ok := subROnlyUsers[bbUser.User.EmailAddress]
		if !ok {
			continue
		}
		subUsers[user.BbEmail] = user
	}

	// "Open" events
	for _, user := range subUsers {
		events = append(events, buildPROpenEvent(repo, pr, user))
	}

	// "Activity" events
	prActivities, err := f.bbClient.ListPRActivity(ctx, repo.Name, pr.ID)
	if err != nil {
		logger.Error().Err(err).Msg("List PR activities")
		return
	}
	for _, activity := range prActivities {
		for _, user := range subUsers {
			if event := buildPRActivityEvent(repo, pr, activity, user); event != nil {
				events = append(events, *event)
			}
		}
	}

	return
}

// getPRRelatedUsers returns the users subscribed to the given repo PRs in two groups: subscribed to all events and subscribed reviewer only events.
func (f *Fetcher) getPRRelatedUsers(ctx context.Context, repoID int64) (subAllUsers, subROnlyUsers userSet, retErr error) {
	_, logger := logging.GetCtxLogger(ctx)

	subAllUsers, subROnlyUsers = make(userSet), make(userSet)

	subs, err := f.storage.ListSubscriptionsForRepoID(ctx, repoID)
	if err != nil {
		return nil, nil, fmt.Errorf("list subscriptions for repo: %w", err)
	}
	for _, sub := range subs {
		user, err := f.storage.GetUserByID(ctx, sub.UserID)
		if err != nil {
			return nil, nil, fmt.Errorf("get user by id (%d): %w", sub.UserID, err)
		}
		if user == nil {
			logger.Warn().
				Int64("sub_id", sub.ID).
				Int64("user_id", sub.UserID).
				Msg("Subscribed user not found")
			continue
		}
		if user.BbEmail == "" {
			logger.Warn().
				Int64("user_id", sub.UserID).
				Msg("Subscribed user has no BitBucket email set")
			continue
		}

		switch sub.Type {
		case model.SubscriptionTypeReviewerOnly:
			subROnlyUsers[user.BbEmail] = *user
		case model.SubscriptionTypeAll:
			subAllUsers[user.BbEmail] = *user
		}
	}

	return subAllUsers, subROnlyUsers, nil
}
