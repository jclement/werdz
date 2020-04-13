package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},           // All origins
		AllowedMethods: []string{"GET", "POST"}, // Allowing only get, just an example
	})

	log.Fatal(http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/test/", a.sayHello).Methods("POST")
	a.Router.HandleFunc("/books/{title}/page/{page}", a.sampleBooks).Methods("GET")
}

func (a *App) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
