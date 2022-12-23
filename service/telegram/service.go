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
	updateCh    tgbotapi.UpdatesChannel
	storage     *storage.Psql
	reconnectCh chan struct{}
}

// NewService creates a new Service instance.
func NewService(ctx context.Context, st *storage.Psql) (*Service, error) {
	if st == nil {
		return nil, fmt.Errorf("storage is nil")
	}

	svc := &Service{
		storage:     st,
		reconnectCh: make(chan struct{}, 1),
	}
	if err := svc.connectToTG(ctx); err != nil {
		return nil, fmt.Errorf("connecting to Telegram: %w", err)
	}

	return svc, nil
}

// Logger returns a logger instance.
func (svc *Service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.KeyService, "telegram").Logger()

	return &logger
}

// connectToTG connects / reconnects to the Telegram.
func (svc *Service) connectToTG(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(viper.GetString(config.TelegramToken))
	if err != nil {
		return fmt.Errorf("bot init: %w", err)
	}

	lastUpdateID, err := svc.getLastUpdateID(ctx)
	if err != nil {
		return fmt.Errorf("reading lastUpdateID: %w", err)
	}

	updateConfig := tgbotapi.NewUpdate(lastUpdateID)
	updateConfig.Timeout = 5
	updateCh := bot.GetUpdatesChan(updateConfig)

	svc.bot, svc.updateCh = bot, updateCh
	svc.Logger(ctx).Info().Str("user", bot.Self.UserName).Msg("Bot initialized")

	return nil
}

// disconnectFromTG disconnects from the Telegram.
// Creates a dummy updateCh to prevent a high CPU usage during the select{}
func (svc *Service) disconnectFromTG() {
	if svc.bot == nil {
		return
	}

	svc.bot.StopReceivingUpdates()
	svc.bot = nil
	svc.updateCh = make(tgbotapi.UpdatesChannel)
}
