package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameHasPlayerRequest struct {
	PlayerID string `json:"playerId"`
}

func (a *App) apiGameHasPlayer(w http.ResponseWriter, r *http.Request) {
	var payload apiGameHasPlayerRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}
	playerID := game.PlayerID(payload.PlayerID)

	vars := mux.Vars(r)
	id := vars["id"]

	if g, gameFound := a.games[game.GID(id)]; gameFound {
		webservice.RespondWithJSON(w, http.StatusOK, g.Game.PlayerExists(playerID))
	} else {
		webservice.RespondWithError(w, http.StatusNotFound, "game does not exist")

	}
}
