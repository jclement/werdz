package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/util/webservice"
)

type apiGameNewRequest struct {
	Rounds int `json:"rounds"`
}

type apiGameNewResponse struct {
	ID string `json:"id"`
}

func (a *App) apiGameNew(w http.ResponseWriter, r *http.Request) {
	var payload apiGameNewRequest
	if err := webservice.HandleJSONRequest(w, r, &payload); err != nil {
		return
	}
	wordFunc := func() (word, definition string) {
		w := a.WordSet.Random()
		return w.Word, w.Definition
	}
	g, _ := game.NewGame(wordFunc, game.ModeNormal, payload.Rounds, 300, 90)

	a.games[g.ID] = &gameState{
		Game:    g,
		Clients: make(map[*websocket.Conn]game.PlayerID),
	}

	webservice.RespondWithJSON(w, http.StatusOK, apiGameNewResponse{ID: string(g.ID)})
}
