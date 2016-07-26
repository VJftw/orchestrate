package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
	"github.com/vjftw/orchestrate/commander/controllers"
	"github.com/vjftw/orchestrate/commander/managers"
	"github.com/vjftw/orchestrate/commander/persisters"
	"github.com/vjftw/orchestrate/commander/providers"
	"github.com/vjftw/orchestrate/commander/resolvers"
	"github.com/vjftw/orchestrate/commander/routers"
	"github.com/vjftw/orchestrate/commander/validators"
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

	var userController controllers.User

	err := orchestrateApp.graph.Provide(
		&inject.Object{Name: "persister.gorm", Value: persisters.NewGORM()},
		&inject.Object{Name: "manager.default", Value: &managers.Model{}},
		&inject.Object{Name: "validator.user", Value: &validators.User{}},
		&inject.Object{Name: "provider.user", Value: providers.NewUser()},
		&inject.Object{Name: "provider.auth_token", Value: providers.NewAuthToken()},
		&inject.Object{Name: "resolver.user", Value: &resolvers.User{}},
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
