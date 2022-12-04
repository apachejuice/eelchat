package config

import (
	"github.com/apachejuice/eelchat/internal/logs"
)

// the logger for this package
var configLogger logs.Logger

func SetupConfig() {
	configLogger = logs.NewLogger("config")

	configLogger.Debug("Loading configuration from file", "filename", "eelchat.json")
	setEnv()
}
