package webservice

import (
	"encoding/json"
	"io"
	"net/http"
)

// RespondWithError responds with a friendly JSON error
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON responds with some JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HandleJSONRequest parses a request and raises an error if necessary
func HandleJSONRequest(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	limitedBodyReader := io.LimitedReader{R: r.Body, N: 1048576}
	decoder := json.NewDecoder(&limitedBodyReader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return err
	}
	defer r.Body.Close()
	return nil
}
