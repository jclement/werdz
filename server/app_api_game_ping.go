package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type gamePingRequest struct {
	PlayerID     string `json:"playerId"`
	RoundID      string `json:"roundId"`
	DefinitionID string `json:"definitionId"`
}

func (a *App) apiGamePing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := game.GID(vars["id"])

	var payload gamePingRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}

	playerID := game.PlayerID(payload.PlayerID)

	if g, ok := a.getGame(id); ok {
		g.LastPing[playerID] = time.Now()
		if inactive, err := g.Game.IsPlayerInactive(playerID); err == nil && inactive {
			g.Game.SetPlayerInactive(playerID, false)
			g.PushUpdate()
		}
	}

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}
