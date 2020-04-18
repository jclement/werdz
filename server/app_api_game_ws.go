package main

import (
	"log"
	"net/http"
	"time"

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
	g, gameFound := a.getGame(game.GID(id))
	if !gameFound {
		webservice.RespondWithError(w, http.StatusNotFound, "game not found")
		return
	}

	playerID := game.PlayerID(r.URL.Query().Get("playerid"))
	if len(playerID) == 0 {
		webservice.RespondWithError(w, http.StatusNotFound, "playerid required")
		return
	}

	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		webservice.RespondWithError(w, http.StatusNotFound, "name required")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		webservice.RespondWithError(w, http.StatusInternalServerError, "can't upgrade")
		return
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	g.Clients[ws] = game.PlayerID(playerID)
	g.LastPing[playerID] = time.Now()
	// join or re-join the game as necessary
	if !g.Game.PlayerExists(playerID) {
		g.Game.AddPlayer(playerID, name)
	} else {
		g.Game.SetPlayerInactive(playerID, false)
	}
	// give them the current state
	m := newGameStateMessage(g.Game, playerID)
	ws.WriteJSON(m)

	// let the world know we have a new player
	g.PushUpdate()
}
