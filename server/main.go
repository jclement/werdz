package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"gitlab.adipose.net/jeff/werdz/util/mattermost"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Default configuration settings
	viper.SetDefault("Listen", "localhost:8100")

	// Parse configuration files
	viper.SetConfigName("werdz.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	a := App{}
	a.Initialize()

	a.SetMattermostWebhook(mattermost.New(
		viper.GetString("WebhookURL"),
		viper.GetString("WebhookUser"),
		viper.GetString("WebhookImageURL"),
	))

	a.Run(viper.GetString("Listen"))
}
