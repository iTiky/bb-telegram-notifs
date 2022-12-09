package config

import (
	"github.com/spf13/viper"
)

const (
	BBPrefix = "BitBucket/"

	BBToken   = BBPrefix + "Token"   // BitBucket API token [string]
	BBHost    = BBPrefix + "Host"    // BitBucket host [string]
	BBProject = BBPrefix + "Project" // BitBucket project [string]
)

func bbSectionSetDefaults() {
	viper.SetDefault(BBToken, "")
	viper.SetDefault(BBHost, "")
	viper.SetDefault(BBProject, "")
}
