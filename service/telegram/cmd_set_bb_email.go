package telegram

import (
	"context"
	"strings"

	"github.com/itiky/bb-telegram-notifs/pkg"
)

// handleSetBBEmailCmd handles BitBucket email set for the user request.
// /set_bb_email cmd.
func (svc *Service) handleSetBBEmailCmd(ctx context.Context, chatID int64, text string) {
	cmdParts := strings.Split(text, " ")
	if len(cmdParts) != 2 {
		svc.sendInvalidFormatMsg(ctx, chatID, cmdSetBBEmailFormat)
		return
	}

	bbEmail := cmdParts[1]
	if err := svc.storage.SetBBEmail(ctx, pkg.GetUserCtx(ctx).ID, bbEmail); err != nil {
		svc.Logger(ctx).Error().Err(err).Msg("Updating BitBucket email")
		svc.sendErrorMsg(ctx, chatID)
		return
	}
	svc.sendMsg(ctx, chatID, "👌 BitBucket email updated")
	svc.Logger(ctx).Info().Msg("BitBucket email updated")
}
