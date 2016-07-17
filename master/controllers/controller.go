package controllers

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Controller - Interface that defines methods that all controllers should have
type Controller interface {
	Setup(*mux.Router, *render.Render)
}
