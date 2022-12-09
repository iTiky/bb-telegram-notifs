package pkg

import (
	"context"

	"github.com/google/uuid"

	"github.com/itiky/bb-telegram-notifs/model"
)

// ContextKey is a context key type.
type ContextKey string

// String implements fmt.Stringer, used to avoid potential context key collisions.
func (c ContextKey) String() string {
	return "bb-" + string(c)
}

const (
	contextKeyUser          = ContextKey("User")          // key to store model.User
	contextKeyCorrelationID = ContextKey("CorrelationID") // key to store unique request ID
)

// ContextWithUser sets the model.User to the context.
func ContextWithUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}

// ContextWithCorrelationID sets the unique request ID to the context.
func ContextWithCorrelationID(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyCorrelationID, uuid.New().String())
}

// GetUserCtx returns the model.User from the context.
func GetUserCtx(ctx context.Context) model.User {
	if ctx == nil {
		return model.User{}
	}

	if ctxValue := ctx.Value(contextKeyUser); ctxValue != nil {
		if user, ok := ctxValue.(model.User); ok {
			return user
		}
	}

	return model.User{}
}

// GetCorrelationIDCtx returns the unique request ID from the context.
func GetCorrelationIDCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if ctxValue := ctx.Value(contextKeyCorrelationID); ctxValue != nil {
		if correlationID, ok := ctxValue.(string); ok {
			return correlationID
		}
	}

	return ""
}
