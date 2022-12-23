package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"

	"github.com/itiky/bb-telegram-notifs/model"
	"github.com/itiky/bb-telegram-notifs/pkg"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
)

const (
	kvLastUpdateID = "tg_last_update_id"
)

// StartUpdatesHandling starts handling Telegram updates.
func (svc *Service) StartUpdatesHandling(ctx context.Context) error {
	go func() {
		logger := svc.Logger(ctx)
		ctx := logging.SetCtxLogger(ctx, *logger)

		for work := true; work; {
			select {
			case <-ctx.Done():
				work = false
			case update, ok := <-svc.updateCh:
				if !ok {
					// Prevent a high CPU usage and wait for reconnectCh to be triggered
					time.Sleep(1 * time.Second)
					break
				}
				svc.handleUpdate(ctx, update)
			case <-svc.reconnectCh:
				logger.Info().Msg("Resubscribing to Telegram updates")
				svc.disconnectFromTG()

				if err := svc.connectToTG(ctx); err != nil {
					logger.Error().
						Err(err).
						Msg("Reconnecting to Telegram")
					time.Sleep(5 * time.Second)
					svc.reconnectCh <- struct{}{}
				}
			}
		}
	}()

	return nil
}

// handleUpdate handles a single Telegram update.
// Method acks the update immediately, check the sender authorization and passes the update to the appropriate handler.
func (svc *Service) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	ctx, logger := logging.SetCtxLoggerStrFields(ctx, "update_id", strconv.Itoa(update.UpdateID))

	// Update last handled update ID (ack update)
	if err := svc.setLastUpdateID(ctx, update.UpdateID); err != nil {
		logger.Fatal().Err(err).Msg("Setting lastUpdateID")
	}

	// Set correlationID
	ctx = pkg.ContextWithCorrelationID(ctx)
	ctx, logger = logging.SetCtxLoggerStrFields(ctx,
		logging.KeyCorrelationID, pkg.GetCorrelationIDCtx(ctx),
	)

	// Get Telegram userName and update type
	var tgLogin string
	var tgID, tgChatID int64

	updIsMessage := false
	if update.Message != nil {
		if !update.Message.IsCommand() {
			logger.Trace().Msg("Message is not a command (skip)")
			return
		}
		if update.Message.From == nil {
			logger.Trace().Msg("Message.From is nil (skip)")
			return
		}
		if update.Message.Chat == nil {
			logger.Trace().Msg("Message.Chat is nil (skip)")
			return
		}

		tgLogin, updIsMessage = update.Message.From.UserName, true
		tgID, tgChatID = update.Message.From.ID, update.Message.Chat.ID
	}

	updIsCallback := false
	if update.CallbackQuery != nil {
		if update.CallbackQuery.From == nil {
			logger.Trace().Msg("CallbackQuery.From is nil (skip)")
			return
		}

		tgLogin, updIsCallback = update.CallbackQuery.From.UserName, true
		tgID = update.CallbackQuery.From.ID
	}

	// Check auth: user must be registered in the system
	ctx, logger, err := svc.enrichCtxForHandling(ctx, tgLogin, tgID, tgChatID)
	if err != nil {
		logger.Info().Err(err).Msg("User not authorized")
		return
	}

	// Handle command message
	if updIsMessage {
		cmd, chatID, msgText := update.Message.Command(), update.Message.Chat.ID, update.Message.Text

		ctx, logger := logging.SetCtxLoggerStrFields(ctx, "command", cmd)
		logger.Info().Msg("Handling command")

		switch cmd {
		case cmdRepos:
			svc.handleReposCmd(ctx, chatID)
		case cmdSetBBEmail:
			svc.handleSetBBEmailCmd(ctx, chatID, msgText)
		default:
			svc.handleHelpCmd(ctx, chatID)
		}
	}

	// Handle callback query
	if updIsCallback {
		cbID, cbData, chatID := update.CallbackQuery.ID, update.CallbackQuery.Data, update.CallbackQuery.Message.Chat.ID

		ctx, logger := logging.SetCtxLoggerStrFields(ctx, "callback", cbData)
		logger.Info().Msg("Handling callback")

		cbDataParts := strings.Split(cbData, "/")
		switch cbDataParts[0] {
		case callbackDataSubscribeAll:
			svc.handleSubscribeCallback(ctx, cbID, chatID, model.SubscriptionTypeAll, cbDataParts[1:]...)
		case callbackDataSubscribeReviewerOnly:
			svc.handleSubscribeCallback(ctx, cbID, chatID, model.SubscriptionTypeReviewerOnly, cbDataParts[1:]...)
		case callbackDataUnsubscribe:
			svc.handleUnsubscribeCallback(ctx, cbID, chatID, cbDataParts[1:]...)
		}
	}
}

// enrichCtxForHandling enriches the context with the user data, correlation ID and logger.
func (svc *Service) enrichCtxForHandling(ctx context.Context, tgLogin string, tgID, tgChatID int64) (context.Context, zerolog.Logger, error) {
	ctx, logger := logging.SetCtxLoggerStrFields(ctx,
		"tg_user", tgLogin,
		"tg_id", strconv.FormatInt(tgID, 10),
		logging.KeyCorrelationID, pkg.GetCorrelationIDCtx(ctx),
	)

	user, err := svc.storage.GetUserByTgID(ctx, tgID)
	if err != nil {
		return ctx, logger, fmt.Errorf("storage.GetUserByTgID: %w", err)
	}
	if user == nil {
		user, err = svc.storage.CreateUser(ctx, tgID, tgChatID, tgLogin)
		if err != nil {
			return ctx, logger, fmt.Errorf("creating inactive user: %w", err)
		}
		logger.Info().Msg("Inactive user created")
	}

	if !user.Active {
		return ctx, logger, fmt.Errorf("user is inactive")
	}

	ctx = pkg.ContextWithUser(ctx, *user)
	ctx = logging.SetCtxLogger(ctx, logger)

	return ctx, logger, nil
}

// getLastUpdateID returns the last handled update ID.
func (svc *Service) getLastUpdateID(ctx context.Context) (int, error) {
	idBz, err := svc.storage.GetKeyValue(ctx, kvLastUpdateID)
	if err != nil {
		return 0, fmt.Errorf("storage.GetKeyValue: %w", err)
	}
	if idBz == nil {
		return 0, nil
	}

	id, err := strconv.ParseInt(*idBz, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt (%s): %w", idBz, err)
	}

	return int(id), nil
}

// setLastUpdateID sets the last handled update ID.
func (svc *Service) setLastUpdateID(ctx context.Context, id int) error {
	if err := svc.storage.SetKeyValue(ctx, kvLastUpdateID, strconv.Itoa(id)); err != nil {
		return fmt.Errorf("storage.SetKeyValue: %w", err)
	}

	return nil
}
