package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Interface that defines methods that all controllers should have
type Controller interface {
	AddRoutes(mux.Router)
}

// Respond - Writes the given status code and object to the response
func Respond(w http.ResponseWriter, code int, v interface{}) {

	r := render.New(render.Options{
		IndentJSON: true,
	})

	r.JSON(w, code, v)
}
