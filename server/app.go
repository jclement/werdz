package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gitlab.adipose.net/jeff/werdz/util/mattermost"
)

// App represents an instance of the application server
type App struct {
	Router     *mux.Router
	Mattermost mattermost.Webhook
}

// Initialize initializes up the application server
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// SetMattermostWebhook sets a webhook for mattermost logging
func (a *App) SetMattermostWebhook(w mattermost.Webhook) {
	a.Mattermost = w
}

// Run starts the application server runnning
func (a *App) Run(addr string) {
	n := negroni.New()

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(a.Router)

	log.Printf("Starting service on %s", addr)
	log.Fatal(http.ListenAndServe(addr, n))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/test/", a.sayHello).Methods("POST")
	a.Router.HandleFunc("/api/player/generate", a.generatePlayerID).Methods("GET")
	a.Router.HandleFunc("/books/{title}/page/{page}", a.sampleBooks).Methods("GET")
}
