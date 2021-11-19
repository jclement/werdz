package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"gitlab.adipose.net/jeff/werdz/models/fakewords"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/models/words"
)

type gameState struct {
	Game          *game.Game
	Clients       map[*websocket.Conn]game.PlayerID
	LastPing      map[game.PlayerID]time.Time
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
	router    *mux.Router
	games     map[game.GID]*gameState
	realWords words.WordSet
	fakeWords fakewords.FakeWordSet
}

// Initialize initializes up the application server
func (a *App) Initialize() {
	a.router = mux.NewRouter()
	a.initializeRoutes()
	a.games = make(map[game.GID]*gameState)
}

// Run starts the application server runnning
func (a *App) Run(addr string, staticDir string) {
	n := negroni.New()

	if staticDir != "" {
		a.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		})

	}

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	if staticDir != "" {
		n.Use(negroni.NewStatic(http.Dir(staticDir)))
	}
	n.UseHandler(a.router)

	// static

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
	a.router.HandleFunc("/api/game/{id}/ping", a.apiGamePing).Methods("POST")
	a.router.HandleFunc("/api/game/count", a.apiGameCount).Methods("GET")
}

// getGame loads from disk if the game has been saved
func (a *App) getGame(id game.GID) (*gameState, bool) {
	if g, ok := a.games[id]; ok {
		return g, true
	}
	return nil, false
}
