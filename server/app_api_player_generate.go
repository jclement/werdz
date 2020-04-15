package main

import (
	"net/http"

	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiPlayerGenerateResponse struct {
	ID string `json:"id"`
}

func (a *App) apiPlayerGenerate(w http.ResponseWriter, r *http.Request) {
	webservice.RespondWithJSON(w, http.StatusOK, apiPlayerGenerateResponse{ID: string(game.GeneratePlayerID())})
}
