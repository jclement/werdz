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
	id := vars["id"]

	g, ok := a.getGame(game.GID(id))
	if !ok {
		webservice.RespondWithError(w, http.StatusNotFound, "unable to find game")
		return
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	if err := g.Game.StartGame(); err != nil {
		webservice.RespondWithError(w, http.StatusInternalServerError, "unable to start game")
		return
	}

	g.PushUpdate()

	a.webhook.Post(fmt.Sprintf("Starting new game %s", id))

	webservice.RespondWithJSON(w, http.StatusOK, nil)
}
