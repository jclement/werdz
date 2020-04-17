package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/models/words"
	"gitlab.adipose.net/jeff/werdz/util/mattermost"
)

type gameState struct {
	Game          *game.Game
	Clients       map[*websocket.Conn]game.PlayerID
	lock          sync.Mutex
	broadcastChan chan bool
}

func (g *gameState) PushUpdate() {
	go func() {
		g.broadcastChan <- true
	}()
}

// App represents an instance of the application server
type App struct {
	router  *mux.Router
	webhook mattermost.Webhook
	games   map[game.GID]*gameState
	WordSet words.WordSet
}

// Initialize initializes up the application server
func (a *App) Initialize() {
	a.router = mux.NewRouter()
	a.initializeRoutes()
	a.games = make(map[game.GID]*gameState)
}

// SetMattermostWebhook sets a webhook for mattermost logging
func (a *App) SetMattermostWebhook(w mattermost.Webhook) {
	a.webhook = w
}

// Run starts the application server runnning
func (a *App) Run(addr string) {
	n := negroni.New()

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(a.router)

	log.Printf("Starting service on %s", addr)

	go a.loop()

	log.Fatal(http.ListenAndServe(addr, n))
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/api/player/generate", a.apiPlayerGenerate).Methods("GET")

	a.router.HandleFunc("/api/game/new", a.apiGameNew).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/start", a.apiGameStart).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/submit", a.apiGameSubmit).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/vote", a.apiGameVote).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/ws", a.apiGameWs)
	a.router.HandleFunc("/api/game/{id}/exists", a.apiGameExists).Methods("GET")
	a.router.HandleFunc("/api/game/{id}/has_player", a.apiGameHasPlayer).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/name_available", a.apiGameNameAvailable).Methods("POST")
}

func (a *App) getGame(id game.GID) (*gameState, bool) {
	if g, ok := a.games[id]; ok {
		return g, true
	}
	return nil, false
}
