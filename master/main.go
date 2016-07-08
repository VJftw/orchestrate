package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
	"github.com/vjftw/orchestrate/master/controllers"
)

// OrchestrateApp - Orchestrate application struct
type OrchestrateApp struct {
	graph  *inject.Graph
	Router *MuxRouter `inject:"default.router"`
}

// NewOrchestrateApp - Initialise with Depencency Injection
func NewOrchestrateApp() *OrchestrateApp {
	orchestrateApp := OrchestrateApp{}
	orchestrateApp.graph = new(inject.Graph)

	// var gormPersister persisters.GORMPersister

	muxRouter := NewMuxRouter()

	orchestrateApp.graph.Provide(
		&inject.Object{Value: &orchestrateApp},
		// &inject.Object{Name: "persister.gorm", Value: &gormPersister},
		&inject.Object{Name: "default.router", Value: muxRouter},
		&inject.Object{
			Name:  "controller.user",
			Value: controllers.NewUserController(muxRouter.Router),
		},
	)

	if err := orchestrateApp.graph.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(orchestrateApp.graph.Objects())

	return &orchestrateApp
}

// AppEngine - For use in Google AppEngine
func AppEngine() {
	app := NewOrchestrateApp()
	http.Handle("/", app.Router.Handler)
}

func main() {
	app := NewOrchestrateApp()
	log.Fatal(http.ListenAndServe(":8734", app.Router.Handler))
}
