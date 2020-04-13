package main

import (
	"fmt"
	"net/http"

	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type helloRequest struct {
	Name string `json:"name"`
}

type helloResponse struct {
	Msg string `json:"msg"`
}

func (a *App) sayHello(w http.ResponseWriter, r *http.Request) {
	var payload helloRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}
	if payload.Name == "Jeff" {
		panic("ahh")
	}
	webservice.RespondWithJSON(w, http.StatusOK, helloResponse{Msg: fmt.Sprintf("Hello %s!", payload.Name)})
}
