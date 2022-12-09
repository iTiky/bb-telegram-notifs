package config

import (
	"github.com/spf13/viper"
)

const (
	TelegramPrefix = "Telegram/"

	TelegramToken = TelegramPrefix + "Token" // Telegram bot token [string]
)

func telegramSectionSetDefaults() {
	viper.SetDefault(TelegramToken, "")
}
