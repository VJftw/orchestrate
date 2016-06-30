package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/persisters"
)

// OrchestrateApp - Orchestrate application struct
type OrchestrateApp struct {
	Router *Router `inject:""`
}

func initApp() OrchestrateApp {
	var g inject.Graph
	var app OrchestrateApp
	var gormDB persisters.GORMPersister

	err := g.Provide(
		&inject.Object{Value: &app},
		&inject.Object{Name: "persister gorm", Value: &gormDB},
		&inject.Object{Name: "manager entity", Value: &managers.EntityManager{}},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gormDB.Init()

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app.Router.init()

	return app
}

// AppEngine - For use in Google AppEngine
func AppEngine() {
	app := initApp()

	http.Handle("/", app.Router.Router)
}

func main() {
	app := initApp()
	log.Fatal(http.ListenAndServe(":8734", app.Router.Router))
}
