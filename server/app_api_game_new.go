package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameNewRequest struct {
	Rounds int `json:"rounds"`
	Mode   int `json:"mode"`
}

type apiGameNewResponse struct {
	ID string `json:"id"`
}

func (a *App) apiGameNew(w http.ResponseWriter, r *http.Request) {
	var payload apiGameNewRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}
	var wordFunc func() (word, definition string)

	if payload.Mode == 0 {
		wordFunc = func() (word, definition string) {
			w := a.realWords.Random()
			return w.Word, w.Definition
		}
	} else if payload.Mode == 1 {
		wordFunc = func() (word, definition string) {
			w := a.fakeWords.Random()
			return w, ""
		}
	} else {
		webservice.RespondWithError(w, http.StatusInternalServerError, "unsupport game mode")
		return
	}

	g, _ := game.NewGame(wordFunc, game.Mode(payload.Mode), payload.Rounds, 60, 30)

	newGame := &gameState{
		Game:          g,
		Clients:       make(map[*websocket.Conn]game.PlayerID),
		broadcastChan: make(chan bool),
		LastPing:      make(map[game.PlayerID]time.Time),
	}

	a.games[g.ID] = newGame
	go a.gameLoop(newGame)

	webservice.RespondWithJSON(w, http.StatusOK, apiGameNewResponse{ID: string(g.ID)})
}
