package fetcher

import (
	"context"

	"github.com/itiky/bb-telegram-notifs/model"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
	bbModel "github.com/itiky/bb-telegram-notifs/provider/bitbucket/model"
)

// handleEvents handles fetched events and passes them through to the Telegram service.
func (f *Fetcher) handleEvents(ctx context.Context, events []model.Event) {
	for _, e := range events {
		ctx, logger := logging.SetCtxLoggerStrFields(ctx,
			"event_hash", e.Hash,
		)

		eventCreated, err := f.storage.CreateEventSafe(ctx, e)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create event")
			continue
		}
		if !eventCreated {
			continue
		}

		logger.Info().Msgf("Sending event: %s", e.String())
		f.tgService.SendEventMessage(ctx, e)
	}
}

// buildPROpenEvent builds a new event for a new PR.
func buildPROpenEvent(repo model.Repo, pr bbModel.PR, recipient model.User) model.Event {
	e := model.Event{
		Type:              model.EventTypePROpen,
		RecipientTgID:     recipient.TgID,
		RecipientTgChatID: recipient.TgChatID,
		SenderName:        pr.Author.User.DisplayName,
		RepoProject:       repo.Project,
		RepoName:          repo.Name,
		PrID:              pr.ID,
		PrTitle:           pr.Title,
		PrURL:             pr.SelfLink().Href,
		CreatedAt:         pr.CreatedAt(),
	}
	e.SetHash()

	return e
}

// buildPRActivityEvent builds a new event for a new PR activity.
func buildPRActivityEvent(repo model.Repo, pr bbModel.PR, activity bbModel.PRActivity, recipient model.User) *model.Event {
	e := model.Event{
		RecipientTgID:     recipient.TgID,
		RecipientTgChatID: recipient.TgChatID,
		SenderName:        activity.User.DisplayName,
		RepoProject:       repo.Project,
		RepoName:          repo.Name,
		PrID:              pr.ID,
		PrTitle:           pr.Title,
		PrURL:             pr.SelfLink().Href,
		CreatedAt:         activity.CreatedAt(),
	}

	switch activity.Action {
	case "APPROVED":
		e.Type = model.EventTypePRApproved
	case "UPDATED":
		e.Type = model.EventTypePRUpdated
	case "COMMENTED":
		e.Type = model.EventTypeComment
	case "MERGED":
		e.Type = model.EventTypePRMerged
	case "DECLINED":
		e.Type = model.EventTypePRRejected
	default:
		return nil
	}
	e.SetHash()

	return &e
}
