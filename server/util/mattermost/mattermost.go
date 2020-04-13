package mattermost

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Webhook represents a mattermost webhook for posting stuff
type Webhook struct {
	url      string
	name     string
	imageURL string
}

// New generates a new Webhook
func New(url string, name string, imageURL string) Webhook {
	return Webhook{
		url:      url,
		imageURL: imageURL,
		name:     name,
	}
}

// Post posts a message to the webhook
func (w Webhook) Post(message string) {
	if w.url == "" {
		return
	}
	reqBody, _ := json.Marshal(map[string]string{
		"username": w.name,
		"icon_url": w.imageURL,
		"text":     message,
	})
	resp, err := http.Post(w.url, "text/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
}
