package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller - Interface that defines methods that all controllers should have
type Controller interface {
	AddRoutes(mux.Router)
}

// Respond - Writes the given status code and object to the response
func Respond(w http.ResponseWriter, code int, v interface{}) {

	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(v)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(code)
	w.Write(j)
}
