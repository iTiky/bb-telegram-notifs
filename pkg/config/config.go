package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Setup initializes the Viper.
func Setup() error {
	viper.SetEnvPrefix("BBTT_")
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer("/", "__"))

	appSectionSetDefaults()
	dbSectionSetDefaults()
	telegramSectionSetDefaults()
	bbSectionSetDefaults()

	if err := appSectionValidate(); err != nil {
		return fmt.Errorf("app: validate: %w", err)
	}

	return nil
}
