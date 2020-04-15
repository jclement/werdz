package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

func (a *App) apiGameStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	id := vars["id"]

	g, ok := a.games[game.GID(id)]
	if !ok {
		webservice.RespondWithError(w, http.StatusInternalServerError, "unable to find game")
		return
	}
	
	if err := g.Game.StartGame(); err != nil {
		webservice.RespondWithError(w, http.StatusInternalServerError, "unable to start game")
		return
	}

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}