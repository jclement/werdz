package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

func post(msg string) {
	if !viper.IsSet("Webhook") {
		return
	}

	url := viper.GetString("Webhook")
	user := viper.GetString("WebhookUser")
	image := viper.GetString("WebhookImage")

	reqBody, _ := json.Marshal(map[string]string{
		"username": user,
		"icon_url": image,
		"text":     msg,
	})
	resp, err := http.Post(url, "text/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
}
