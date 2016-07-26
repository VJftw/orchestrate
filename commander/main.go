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
	var authController controllers.Auth
	var projectController controllers.Project

	err := orchestrateApp.graph.Provide(
		&inject.Object{Name: "persister.gorm", Value: persisters.NewGORM()},
		&inject.Object{Name: "manager.default", Value: managers.NewModel()},
		&inject.Object{Name: "validator.user", Value: validators.NewUser()},
		&inject.Object{Name: "validator.project", Value: validators.NewProject()},
		&inject.Object{Name: "provider.user", Value: providers.NewUser()},
		&inject.Object{Name: "provider.project", Value: providers.NewProject()},
		&inject.Object{Name: "provider.auth_token", Value: providers.NewAuthToken()},
		&inject.Object{Name: "resolver.user", Value: resolvers.NewUser()},
		&inject.Object{Name: "resolver.project", Value: resolvers.NewProject()},
		&inject.Object{
			Name:  "controller.user",
			Value: &userController,
		},
		&inject.Object{
			Name:  "controller.auth",
			Value: &authController,
		},
		&inject.Object{
			Name:  "controller.project",
			Value: &projectController,
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
		&authController,
		&projectController,
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
