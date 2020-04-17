package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

func (a *App) apiGameExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, gameFound := a.games[game.GID(id)]
	webservice.RespondWithJSON(w, http.StatusOK, gameFound)
}
