package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type helloRequest struct {
	Name string
}

type helloResponse struct {
	Msg string
}

func (a *App) sayHello(w http.ResponseWriter, r *http.Request) {

	var hR helloRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&hR); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	a.respondWithJSON(w, 200, helloResponse{Msg: fmt.Sprintf("Hello %s!", hR.Name)})
}
