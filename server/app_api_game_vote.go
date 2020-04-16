package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameVoteRequest struct {
	PlayerID     string `json:"playerId"`
	RoundID      string `json:"roundId"`
	DefinitionID string `json:"definitionId"`
}

func (a *App) apiGameVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := game.GID(vars["id"])

	var payload apiGameVoteRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}

	playerID := game.PlayerID(payload.PlayerID)
	roundID := game.RoundID(payload.RoundID)
	definitionID := game.DefinitionID(payload.DefinitionID)

	g, ok := a.games[id]
	if !ok {
		webservice.RespondWithError(w, http.StatusNotFound, "game does not exist")
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	if err := g.Game.Vote(playerID, roundID, definitionID); err != nil {
		webservice.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	g.PushUpdate()

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}
