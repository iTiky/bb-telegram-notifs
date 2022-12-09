package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/itiky/bb-telegram-notifs/model"
)

// ListSubscriptionsForUser lists all model.Subscription for the given user.
func (st *Psql) ListSubscriptionsForUser(ctx context.Context, userID int64) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := st.db.NewSelect().Model(&subs).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("selecting subscriptions: %w", err)
	}

	return subs, nil
}

// ListSubscriptionsForRepoID lists all model.Subscription for the given repository.
func (st *Psql) ListSubscriptionsForRepoID(ctx context.Context, repoID int64) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := st.db.NewSelect().Model(&subs).Where("repo_id = ?", repoID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("selecting subscriptions: %w", err)
	}

	return subs, nil
}

// SetSubscription updates the model.Subscription for the given user and repository.
func (st *Psql) SetSubscription(ctx context.Context, userID, repoID int64, subType model.SubscriptionType) error {
	now := time.Now().UTC()
	sub := model.Subscription{
		UserID:    userID,
		RepoID:    repoID,
		Type:      subType,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := st.db.NewInsert().
		Model(&sub).
		On("CONFLICT (user_id, repo_id) DO UPDATE").
		Set("type = ?", sub.Type).
		Set("updated_at = ?", sub.UpdatedAt).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("upserting subscription: %w", err)
	}

	return nil
}

// DeleteSubscription deletes the model.Subscription for the given user and repository.
func (st *Psql) DeleteSubscription(ctx context.Context, userID, repoID int64) error {
	_, err := st.db.NewDelete().Model(&model.Subscription{}).Where("user_id = ? AND repo_id = ?", userID, repoID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting subscription: %w", err)
	}

	return nil
}
