package config

import (
	"github.com/spf13/viper"
)

const (
	DBPrefix = "DB/"

	DBHost     = DBPrefix + "Host"     // PSQL host [string]
	DBPort     = DBPrefix + "Port"     // PSQL port [int]
	DBUser     = DBPrefix + "User"     // PSQL user [string]
	DBPassword = DBPrefix + "Password" // PSQL password [string]
	DBName     = DBPrefix + "Name"     // PSQL database name [string]
	DBSSLMode  = DBPrefix + "SSLMode"  // PSQL connection SSL mode [string, enable / disable]
)

func dbSectionSetDefaults() {
	viper.SetDefault(DBHost, "localhost")
	viper.SetDefault(DBPort, 5432)
	viper.SetDefault(DBUser, "postgres")
	viper.SetDefault(DBPassword, "postgres")
	viper.SetDefault(DBName, "bb-telegram-notifs")
	viper.SetDefault(DBSSLMode, "disable")
}
