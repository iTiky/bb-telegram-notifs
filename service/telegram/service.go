package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/itiky/bb-telegram-notifs/pkg/config"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
	"github.com/itiky/bb-telegram-notifs/provider/storage"
)

// Service is a Telegram bot service that reacts to cmds, callbacks and message send requests.
type Service struct {
	bot         *tgbotapi.BotAPI
	storage     *storage.Psql
	reconnectCh chan struct{}
}

// NewService creates a new Service instance.
func NewService(ctx context.Context, st *storage.Psql) (*Service, error) {
	if st == nil {
		return nil, fmt.Errorf("storage is nil")
	}

	bot, err := tgbotapi.NewBotAPI(viper.GetString(config.TelegramToken))
	if err != nil {
		return nil, fmt.Errorf("bot init: %w", err)
	}

	svc := &Service{
		bot:         bot,
		storage:     st,
		reconnectCh: make(chan struct{}, 1),
	}
	svc.Logger(ctx).Info().Str("user", bot.Self.UserName).Msg("Bot initialized")

	return svc, nil
}

// Logger returns a logger instance.
func (svc *Service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.KeyService, "telegram").Logger()

	return &logger
}
