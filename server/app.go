package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// App represents an instance of the application server
type App struct {
	Router *mux.Router
}

// Initialize initializes up the application server
func (a *App) Initialize() {

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

// Run starts the application server runnning
func (a *App) Run(addr string) {
	n := negroni.New()

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(a.Router)

	log.Fatal(http.ListenAndServe(addr, n))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/test/", a.sayHello).Methods("POST")
	a.Router.HandleFunc("/books/{title}/page/{page}", a.sampleBooks).Methods("GET")
}
