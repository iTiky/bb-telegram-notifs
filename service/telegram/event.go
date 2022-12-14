package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/itiky/bb-telegram-notifs/model"
)

// SendEventMessage sends a message to the user about the BitBucket PR's event.
func (svc *Service) SendEventMessage(ctx context.Context, e model.Event) bool {
	var eventType string
	switch e.Type {
	case model.EventTypePROpen:
		eventType = "📖 PR OPENED"
	case model.EventTypePRApproved:
		eventType = "✅ PR APPROVED"
	case model.EventTypePRRejected:
		eventType = "◀️ PR MERGED"
	case model.EventTypePRMerged:
		eventType = "❌ PR DECLINED"
	case model.EventTypeComment:
		eventType = "💬 COMMENTED"
	case model.EventTypePRUpdated:
		eventType = "👥 PR UPDATED (reviewers changed)"
	}

	text := fmt.Sprintf("*Project:* %s/%s\n*PR:* %s \\[%d]\n*%s*: %s\n*Time:* %s",
		e.RepoProject, e.RepoName,
		e.PrTitle, e.PrID,
		e.SenderName, eventType,
		e.CreatedAt.String(),
	)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Open PR", e.PrURL),
		),
	)

	return svc.sendMsg(ctx, e.RecipientTgChatID, text, withReplyMarkup(markup))
}
