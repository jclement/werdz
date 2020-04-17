package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameNameAvailableRequest struct {
	Name string `json:"name"`
}

func (a *App) apiGameNameAvailable(w http.ResponseWriter, r *http.Request) {
	var payload apiGameNameAvailableRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if g, gameFound := a.games[game.GID(id)]; gameFound {
		webservice.RespondWithJSON(w, http.StatusOK, g.Game.NameAvailable(payload.Name))
	} else {
		webservice.RespondWithError(w, http.StatusNotFound, "game does not exist")

	}
}
