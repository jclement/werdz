package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	// Default configuration settings
	viper.SetDefault("Listen", "localhost:8100")

	// Parse configuration files
	viper.SetConfigName("werdz")
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
	a.Run(viper.GetString("Listen"))
}
