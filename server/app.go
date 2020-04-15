package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"gitlab.adipose.net/jeff/werdz/models/game"
	"gitlab.adipose.net/jeff/werdz/models/words"
	"gitlab.adipose.net/jeff/werdz/util/mattermost"
)

type gameState struct {
	Game    *game.Game
	Clients map[*websocket.Conn]game.PlayerID
}

// App represents an instance of the application server
type App struct {
	router  *mux.Router
	webhook mattermost.Webhook
	games   map[game.GID]gameState
	WordSet words.WordSet
}

// Initialize initializes up the application server
func (a *App) Initialize() {
	a.router = mux.NewRouter()
	a.initializeRoutes()
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

func (a *App) wsHandler(w http.ResponseWriter, r *http.Request) {
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
	player := vars["player"]

	if g, ok := a.games[game.GID(id)]; ok {
		g.Clients[ws] = game.PlayerID(player)
	}

}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/api/player/generate", a.apiPlayerGenerate).Methods("GET")

	a.router.HandleFunc("/api/game/new", a.apiGameNew).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/start", a.apiGameStart).Methods("POST")
	a.router.HandleFunc("/api/game/{id}/{player}/ws", a.wsHandler)
}

func (a *App) loop() {
	for {
		for _, g := range a.games {
			for c := range g.Clients {
				err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					c.Close()
					delete(g.Clients, c)
				}
			}
			// if g.clients == empty kill it and save to disk
		}
		time.Sleep(time.Second)
	}
}
