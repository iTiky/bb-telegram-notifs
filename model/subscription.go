package model

import (
	"time"

	"github.com/uptrace/bun"
)

// SubscriptionType defines how a user would like to receive notifications.
type SubscriptionType string

const (
	SubscriptionTypeAll          SubscriptionType = "all"           // receive all notifications for all PRs
	SubscriptionTypeReviewerOnly SubscriptionType = "reviewer_only" // receive notifications only for PRs where the user is a reviewer
)

// Subscription defines a User's subscription to a repository PR updates.
type Subscription struct {
	bun.BaseModel `bun:"table:subscriptions,alias:s"`

	ID        int64 `bun:",pk,autoincrement"`
	UserID    int64
	RepoID    int64
	Type      SubscriptionType
	CreatedAt time.Time
	UpdatedAt time.Time
}
