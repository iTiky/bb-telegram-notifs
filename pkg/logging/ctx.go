package logging

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/itiky/bb-telegram-notifs/pkg"
	"github.com/itiky/bb-telegram-notifs/pkg/config"
)

// contextKeyLogger is a context key for logger.
const contextKeyLogger = pkg.ContextKey("Logger")

// GetCtxLogger returns a logger from the context or creates a new one and wraps it with the context.
func GetCtxLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}

	if ctxValue := ctx.Value(contextKeyLogger); ctxValue != nil {
		if ctxLogger, ok := ctxValue.(zerolog.Logger); ok {
			return ctx, ctxLogger
		}
	}

	logLevel, _ := zerolog.ParseLevel(viper.GetString(config.LogLevel))
	logger := NewLogger(
		WithLogLevel(logLevel),
	)

	return SetCtxLogger(ctx, logger), logger
}

// SetCtxLogger sets a logger to the context.
func SetCtxLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// SetCtxLoggerStrFields sets a logger to the context with additional string fields.
func SetCtxLoggerStrFields(ctx context.Context, keyvals ...string) (context.Context, zerolog.Logger) {
	ctx, logger := GetCtxLogger(ctx)

	loggerCtx := logger.With()
	for i := 0; i < len(keyvals)-1; i += 2 {
		loggerCtx = loggerCtx.Str(keyvals[i], keyvals[i+1])
	}
	logger = loggerCtx.Logger()

	return SetCtxLogger(ctx, logger), logger
}
