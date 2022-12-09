package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/itiky/bb-telegram-notifs/model"
	"github.com/itiky/bb-telegram-notifs/pkg"
)

// handleReposCmd handles repository list request.
// /repos cmd.
func (svc *Service) handleReposCmd(ctx context.Context, chatID int64) {
	user := pkg.GetUserCtx(ctx)

	repos, err := svc.storage.ListRepos(ctx)
	if err != nil {
		svc.Logger(ctx).Error().Err(err).Msg("Getting repos")
		svc.sendErrorMsg(ctx, chatID)
		return
	}

	subs, err := svc.storage.ListSubscriptionsForUser(ctx, user.ID)
	if err != nil {
		svc.Logger(ctx).Error().Err(err).Msg("Getting subs")
		svc.sendErrorMsg(ctx, chatID)
		return
	}

	for _, repo := range repos {
		var subActive bool
		var subType model.SubscriptionType
		var keyboardBtnsRow []tgbotapi.InlineKeyboardButton

		addSubToAllBtn := func(repoID int64) {
			cbData := fmt.Sprintf("%s/%d/%d", callbackDataSubscribeAll, user.ID, repoID)
			keyboardBtnsRow = append(keyboardBtnsRow, tgbotapi.NewInlineKeyboardButtonData("✅ Subscribe (all)", cbData))
		}
		addSubToReviewerOnlyBtn := func(repoID int64) {
			cbData := fmt.Sprintf("%s/%d/%d", callbackDataSubscribeReviewerOnly, user.ID, repoID)
			keyboardBtnsRow = append(keyboardBtnsRow, tgbotapi.NewInlineKeyboardButtonData("✅ Subscribe (reviewer)", cbData))
		}
		addUnsubBtn := func(repoID int64) {
			cbData := fmt.Sprintf("%s/%d/%d", callbackDataUnsubscribe, user.ID, repoID)
			keyboardBtnsRow = append(keyboardBtnsRow, tgbotapi.NewInlineKeyboardButtonData("❌ Unsubscribe", cbData))
		}

		for _, sub := range subs {
			if sub.RepoID == repo.ID {
				subActive, subType = true, sub.Type
				break
			}
		}

		if !subActive {
			addSubToAllBtn(repo.ID)
			addSubToReviewerOnlyBtn(repo.ID)
		} else {
			addUnsubBtn(repo.ID)
			switch subType {
			case model.SubscriptionTypeAll:
				addSubToReviewerOnlyBtn(repo.ID)
			case model.SubscriptionTypeReviewerOnly:
				addSubToAllBtn(repo.ID)
			}
		}

		markup := tgbotapi.NewInlineKeyboardMarkup(keyboardBtnsRow)
		svc.sendMsg(ctx, chatID, "➡️ "+repo.String(), withReplyMarkup(markup))
	}
}
