package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/sebest/xff"
)

// Router - The application router
type Router struct {
	Router http.Handler
}

func (r *Router) init() {
	router := mux.NewRouter()

	corsHandler := cors.Default()
	xffmw, _ := xff.Default()
	chain := alice.New(corsHandler.Handler, xffmw.Handler)

	r.Router = chain.Then(router)
}
