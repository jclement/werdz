package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/viper"
	"gitlab.adipose.net/jeff/werdz/util/mattermost"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Default configuration settings
	viper.SetDefault("Listen", "localhost:8100")
	viper.SetDefault("Words", "./data/words.json")
	viper.SetDefault("FakeWords", "./data/fake-words.json")

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

	// Load the word list
	rdr, err := os.Open(viper.GetString("Words"))
	if err != nil {
		panic(err)
	}
	defer rdr.Close()
	a.realWords.Load(rdr)

	rdr, err = os.Open(viper.GetString("FakeWords"))
	if err != nil {
		panic(err)
	}
	defer rdr.Close()
	a.fakeWords.Load(rdr)

	a.SetMattermostWebhook(mattermost.New(
		viper.GetString("WebhookURL"),
		viper.GetString("WebhookUser"),
		viper.GetString("WebhookImageURL"),
	))

	a.Run(viper.GetString("Listen"))
}
