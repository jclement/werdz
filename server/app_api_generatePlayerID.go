package main

import (
	"net/http"

	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type generatePlayerIDResponse struct {
	ID string `json:"id"`
}

func (a *App) generatePlayerID(w http.ResponseWriter, r *http.Request) {
	webservice.RespondWithJSON(w, http.StatusOK, generatePlayerIDResponse{ID: string(game.GeneratePlayerID())})
}
