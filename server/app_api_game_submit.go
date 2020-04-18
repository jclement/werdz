package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameSubmitRequest struct {
	PlayerID   string `json:"playerId"`
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

	g, ok := a.getGame(id)
	if !ok {
		webservice.RespondWithError(w, http.StatusNotFound, "game does not exist")
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	if err := g.Game.SubmitWord(playerID, roundID, payload.Definition); err != nil {
		webservice.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		if p, err := g.Game.FindPlayer(playerID); err == nil {
			a.webhook.Post(fmt.Sprintf("[%s / Round %d] Player '%s' submitted '%s' for '%s'", id, len(g.Game.Rounds), p.Name, payload.Definition, g.Game.CurrentRound().Word))
		}
	}

	g.PushUpdate()

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}
