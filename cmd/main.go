package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/itiky/bb-telegram-notifs/pkg/config"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
	"github.com/itiky/bb-telegram-notifs/provider/bitbucket"
	"github.com/itiky/bb-telegram-notifs/provider/storage"
	"github.com/itiky/bb-telegram-notifs/service/fetcher"
	"github.com/itiky/bb-telegram-notifs/service/telegram"
)

func main() {
	logger := logging.NewLogger()

	if err := config.Setup(); err != nil {
		logger.Fatal().Err(err).Msg("Config: init")
	}

	// Init all services with timeout
	initCtx, initCtxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer initCtxCancel()

	st, err := storage.NewPsql(initCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Storage: init")
	}

	tgSvc, err := telegram.NewService(initCtx, st)
	if err != nil {
		logger.Fatal().Err(err).Msg("Telegram: service init")
	}

	bbClient, err := bitbucket.NewClient(initCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("BitBucket: client init")
	}

	fetcherSvc, err := fetcher.NewFetcher(initCtx, bbClient, tgSvc, st)
	if err != nil {
		logger.Fatal().Err(err).Msg("Fetcher: service init")
	}

	// Working ctx with signal notification
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	if err := tgSvc.StartUpdatesHandling(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Telegram: updates handling init")
	}

	if err := fetcherSvc.StartWorkers(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Fetcher: workers init")
	}

	<-ctx.Done()
	logger.Info().Msg("Shutdown")
	time.Sleep(1 * time.Second)
}
