package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.adipose.net/jeff/werdz/models/game"
)

func (a *App) apiGameWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	ws.WriteJSON(struct{ Name string }{
		Name: "Jeff",
	})

	vars := mux.Vars(r)
	id := vars["id"]

	playerID := game.PlayerID(r.URL.Query().Get("playerid"))
	name := r.URL.Query().Get("name")

	if g, ok := a.games[game.GID(id)]; ok {
		g.Lock.Lock()
		defer g.Lock.Unlock()
		g.Clients[ws] = game.PlayerID(id)
		// join or re-join the game as necessary
		if !g.Game.PlayerExists(playerID) {
			g.Game.AddPlayer(playerID, name)
		} else {
			g.Game.SetPlayerInactive(playerID, true)
		}
		// give them the current state
		m := newGameStateMessage(g.Game)
		ws.WriteJSON(m)
		// set game to dirty so everyone else gets an update
		g.Dirty = true
	}

}
