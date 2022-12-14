package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/itiky/bb-telegram-notifs/pkg"
)

// sendMessageOpt is a sendMessage option.
type sendMessageOpt func(cfg *tgbotapi.MessageConfig)

// withReplyMarkup sets a reply markup.
func withReplyMarkup(markup interface{}) sendMessageOpt {
	return func(cfg *tgbotapi.MessageConfig) {
		cfg.ReplyMarkup = markup
	}
}

// sendMsg sends a message to the user.
func (svc *Service) sendMsg(ctx context.Context, chatID int64, text string, opts ...sendMessageOpt) bool {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	for _, opt := range opts {
		opt(&msg)
	}

	if _, err := svc.bot.Send(msg); err != nil {
		svc.Logger(ctx).Error().Err(err).Msg("Sending message")
		return false
	}

	return true
}

// sendErrorMsg sends an error message to the user (correlation ID is only visible to a user).
func (svc *Service) sendErrorMsg(ctx context.Context, chatID int64, opts ...sendMessageOpt) bool {
	text := fmt.Sprintf("❌ Smth went wrong, your correlation ID: %s", pkg.GetCorrelationIDCtx(ctx))
	return svc.sendMsg(ctx, chatID, text, opts...)
}

// sendInvalidFormatMsg sends an invalid format message to the user.
func (svc *Service) sendInvalidFormatMsg(ctx context.Context, chatID int64, format string, opts ...sendMessageOpt) bool {
	text := "🤷‍♂️ Invalid command format, expected:\n" + format
	return svc.sendMsg(ctx, chatID, text, opts...)
}

// sendCallbackMsg sends a callback message to the user.
func (svc *Service) sendCallbackMsg(ctx context.Context, callbackID string, text string) bool {
	callback := tgbotapi.NewCallback(callbackID, text)
	if _, err := svc.bot.Request(callback); err != nil {
		svc.Logger(ctx).Error().Err(err).Msg("Sending callback message")
		return false
	}

	return true
}
