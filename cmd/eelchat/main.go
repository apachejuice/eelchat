package main

import (
	"log"
	"os"

	"github.com/apachejuice/eelchat/internal/api/rest"
	"github.com/apachejuice/eelchat/internal/api/spec"
	"github.com/apachejuice/eelchat/internal/config"
	"github.com/apachejuice/eelchat/internal/db/model"
	"github.com/apachejuice/eelchat/internal/db/repository"
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

	u := repository.NewUserRepository()
	u.Insert(&model.User{})
	u.Count()

	api := rest.NewAPI(
		func(ctx rest.CreateUserContext, rctx rest.RequestContext, user spec.User) spec.CreateUserRes {
			return ctx.InternalServerError(ctx.Error("This is not the way", "Not at all anything"))
		},
	)

	api.ConfigureTLS()
	api.Run()
}
