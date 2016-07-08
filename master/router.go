package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// MuxRouter - The application router
type MuxRouter struct {
	Router  *mux.Router
	Handler http.Handler
}

func NewMuxRouter() *MuxRouter {
	muxRouter := MuxRouter{}

	muxRouter.Router = mux.NewRouter()

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())

	n.UseHandler(muxRouter.Router)

	muxRouter.Handler = n

	return &muxRouter
}
