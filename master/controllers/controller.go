package controllers

import "github.com/gorilla/mux"

// Controller - Interface that defines methods that all controllers should have
type Controller interface {
	AddRoutes(mux.Router)
}
