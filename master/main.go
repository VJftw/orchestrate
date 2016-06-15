package main

import (
	"log"
	"net/http"

	"github.com/facebookgo/inject"
)

// OrchestrateApp - Orchestrate application struct
type OrchestrateApp struct {
	Router *Router `inject:""`
}

func (app *OrchestrateApp) init() {
	inject.Populate(&app)

	app.Router.init()
}

// AppEngine - For use in Google AppEngine
func AppEngine() {
	app := OrchestrateApp{}
	app.init()

	http.Handle("/", app.Router.Router)

}

func main() {
	app := OrchestrateApp{}
	app.init()

	log.Fatal(http.ListenAndServe(":8080", app.Router.Router))
}
