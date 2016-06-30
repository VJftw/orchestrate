package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/vjftw/orchestrate/master/controllers"
)

// Router - The application router
type Router struct {
	Router         http.Handler
	UserController controllers.UserController `inject:"inline"`
	AuthController controllers.AuthController `inject:"inline"`
}

func (r *Router) init() {
	router := mux.NewRouter()

	r.UserController.AddRoutes(router)
	r.AuthController.AddRoutes(router)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(router)

	r.Router = n
}
