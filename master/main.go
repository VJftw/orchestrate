package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
	"github.com/vjftw/orchestrate/master/controllers"
	"github.com/vjftw/orchestrate/master/managers"
	"github.com/vjftw/orchestrate/master/persisters"
	"github.com/vjftw/orchestrate/master/providers"
	"github.com/vjftw/orchestrate/master/resolvers"
	"github.com/vjftw/orchestrate/master/routers"
	"github.com/vjftw/orchestrate/master/validators"
)

// OrchestrateApp - Orchestrate application struct
type OrchestrateApp struct {
	graph  *inject.Graph
	Router *routers.MuxRouter
}

type injectLogger struct {
}

func (l injectLogger) Debugf(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v))
}

// NewOrchestrateApp - Initialise with Depencency Injection
func NewOrchestrateApp() *OrchestrateApp {
	orchestrateApp := OrchestrateApp{}
	orchestrateApp.graph = &inject.Graph{
	// Logger: injectLogger{},
	}

	var userController controllers.UserController

	err := orchestrateApp.graph.Provide(
		&inject.Object{Name: "persister.gorm", Value: persisters.NewGORMPersister()},
		&inject.Object{Name: "manager.default", Value: &managers.ModelManager{}},
		&inject.Object{Name: "validator.user", Value: &validators.UserValidator{}},
		&inject.Object{Name: "provider.user", Value: providers.NewUserProvider()},
		&inject.Object{Name: "resolver.user", Value: &resolvers.UserResolver{}},
		&inject.Object{
			Name:  "controller.user",
			Value: &userController,
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := orchestrateApp.graph.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// for _, element := range orchestrateApp.graph.Objects() {
	// 	fmt.Println(element.Name, &element.Value)
	// }

	muxRouter := routers.NewMuxRouter([]routers.Routable{
		&userController,
	}, true)

	orchestrateApp.Router = muxRouter

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
