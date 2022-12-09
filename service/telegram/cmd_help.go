package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"

	"github.com/itiky/bb-telegram-notifs/pkg"
	"github.com/itiky/bb-telegram-notifs/pkg/config"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
)

// handleHelpCmd sends a welcome message to the user.
// /start cmd.
func (svc *Service) handleHelpCmd(ctx context.Context, chatID int64) {
	text := fmt.Sprintf("Hello, @%s!\nThe current BitBucket project is: %s.",
		pkg.GetUserCtx(ctx).TgLogin,
		viper.GetString(config.BBProject),
	)
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := svc.bot.Send(msg); err != nil {
		_, logger := logging.GetCtxLogger(ctx)
		logger.Error().Err(err).Msg("Sending message")
	}
}
