package controllers

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// IController - Interface that defines methods that all controllers should have
type IController interface {
	Setup(*mux.Router, *render.Render)
}
