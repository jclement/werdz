package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

func (a *App) apiGameWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	vars := mux.Vars(r)
	id := vars["id"]
	g, gameFound := a.games[game.GID(id)]
	if !gameFound {
		webservice.RespondWithError(w, http.StatusNotFound, "game not found")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		webservice.RespondWithError(w, http.StatusInternalServerError, "can't upgrade")
		return
	}

	playerID := game.PlayerID(r.URL.Query().Get("playerid"))
	name := r.URL.Query().Get("name")

	g.Lock.Lock()
	defer g.Lock.Unlock()
	g.Clients[ws] = game.PlayerID(playerID)
	// join or re-join the game as necessary
	if !g.Game.PlayerExists(playerID) {
		g.Game.AddPlayer(playerID, name)
	} else {
		g.Game.SetPlayerInactive(playerID, false)
	}
	// give them the current state
	m := newGameStateMessage(g.Game, playerID)
	ws.WriteJSON(m)
	// set game to dirty so everyone else gets an update
	g.Dirty = true
}
