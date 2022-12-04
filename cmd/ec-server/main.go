package main

import (
	"log"
	"os"

	"github.com/apachejuice/eelchat/internal/api/impl"
	"github.com/apachejuice/eelchat/internal/config"
	"github.com/spf13/viper"
)

func clearLogFile() {
	if err := os.Remove(viper.GetString("logging.file")); err != nil {
		log.Fatal(err)
	}
}

func loadInitialConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/eelchat")

	viper.SetConfigName("eelchat")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadInitialConfig()
	clearLogFile()
	config.SetupConfig()

	impl := impl.NewImpl()
	api := impl.Connect()

	api.ConfigureTLS()
	api.Run()
}
