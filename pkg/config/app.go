package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	LogLevel         = "LogLevel"         // logging level [string]
	FetchPeriod      = "FetchPeriod"      // period between BitBucket API fetches [time.Duration]
	EventsGCPeriod   = "EventsGCPeriod"   // period between events garbage collection (by createdAt) [time.Duration]
	EventGCThreshold = "EventGCThreshold" // events older than this threshold will be deleted [time.Duration]
)

func appSectionSetDefaults() {
	viper.SetDefault(LogLevel, "debug")
	viper.SetDefault(FetchPeriod, 1*time.Minute)
	viper.SetDefault(EventsGCPeriod, 3*time.Hour)
	viper.SetDefault(EventGCThreshold, 7*24*time.Hour)
}

func appSectionValidate() error {
	if _, err := zerolog.ParseLevel(viper.GetString(LogLevel)); err != nil {
		return fmt.Errorf("%s: invalid", LogLevel)
	}

	return nil
}
