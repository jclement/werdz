package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameSubmitRequest struct {
	PlayerID   string `json:"roundId"`
	RoundID    string `json:"roundId"`
	Definition string `json:"definition"`
}

func (a *App) apiGameSubmit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := game.GID(vars["id"])

	var payload apiGameSubmitRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}

	playerID := game.PlayerID(payload.PlayerID)
	roundID := game.RoundID(payload.RoundID)

	g, ok := a.games[id]
	if !ok {
		webservice.RespondWithError(w, http.StatusNotFound, "game does not exist")
	}

	g.Lock.Lock()
	defer g.Lock.Unlock()
	if err := g.Game.SubmitWord(playerID, roundID, payload.Definition); err != nil {
		webservice.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	g.Dirty =  true

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}
