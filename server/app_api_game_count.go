package main

import (
	"net/http"

	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

func (a *App) apiGameCount(w http.ResponseWriter, r *http.Request) {
	cnt := 0
	for _, g := range a.games {
		if g.Game.State == game.StateActive {
			cnt++
		}
	}
	webservice.RespondWithJSON(w, http.StatusOK, cnt)
}
