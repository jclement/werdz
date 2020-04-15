package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

func (a *App) apiGameStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	g, ok := a.games[game.GID(id)]
	if !ok {
		webservice.RespondWithError(w, http.StatusNotFound, "unable to find game")
		return
	}

	g.Lock.Lock()
	defer g.Lock.Unlock()
	if err := g.Game.StartGame(); err != nil {
		webservice.RespondWithError(w, http.StatusInternalServerError, "unable to start game")
		return
	}
	g.Dirty =  true

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}