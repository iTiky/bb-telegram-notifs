package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/itiky/bb-telegram-notifs/model"
)

// handleSubscribeCallback handles subscription callback.
// /repos cmd.
func (svc *Service) handleSubscribeCallback(ctx context.Context, cbID string, chatID int64, subType model.SubscriptionType, opts ...string) {
	userID, repoID, err := parseSubscriptionCallbackData(opts)
	if err != nil {
		svc.Logger(ctx).Error().Strs("opts", opts).Msg("Invalid callback data")
		svc.sendErrorMsg(ctx, chatID)
		return
	}

	svc.sendCallbackMsg(ctx, cbID, "Updating subscription...")

	if err := svc.storage.SetSubscription(ctx, userID, repoID, subType); err != nil {
		svc.Logger(ctx).Error().
			Err(err).
			Int64("userID", userID).
			Int64("repoID", repoID).
			Msg("Updating subscription")
		svc.sendErrorMsg(ctx, chatID)
		return
	}

	var comment string
	switch subType {
	case model.SubscriptionTypeAll:
		comment = "to all events"
	case model.SubscriptionTypeReviewerOnly:
		comment = "only for PR I'm reviewing"
	}
	svc.sendMsg(ctx, chatID, "👌 Subscription updated: "+comment)
	svc.Logger(ctx).Info().
		Int64("userID", userID).
		Int64("repoID", repoID).
		Str("subType", string(subType)).
		Msg("Subscription updated")
}

// handleUnsubscribeCallback handles unsubscription callback.
// /repos cmd.
func (svc *Service) handleUnsubscribeCallback(ctx context.Context, cbID string, chatID int64, opts ...string) {
	userID, repoID, err := parseSubscriptionCallbackData(opts)
	if err != nil {
		svc.Logger(ctx).Error().Msg("Invalid callback data")
		svc.sendErrorMsg(ctx, chatID)
		return
	}

	svc.sendCallbackMsg(ctx, cbID, "Removing subscription...")

	if err := svc.storage.DeleteSubscription(ctx, userID, repoID); err != nil {
		svc.Logger(ctx).Error().Err(err).Msg("Deleting subscription")
		svc.sendErrorMsg(ctx, chatID)
		return
	}

	svc.sendMsg(ctx, chatID, "👌 Subscription removed")
	svc.Logger(ctx).Info().
		Int64("userID", userID).
		Int64("repoID", repoID).
		Msg("Subscription removed")
}

// parseSubscriptionCallbackData splits the callback data into userID and repoID.
func parseSubscriptionCallbackData(data []string) (int64, int64, error) {
	if len(data) != 2 {
		return 0, 0, fmt.Errorf("invalid length")
	}

	userID, err := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing userID: %w", err)
	}

	repoID, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing repoID: %w", err)
	}

	return userID, repoID, nil
}
