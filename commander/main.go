package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/facebookgo/inject"
	"github.com/vjftw/orchestrate/commander/domain/auth"
	"github.com/vjftw/orchestrate/commander/domain/project"
	"github.com/vjftw/orchestrate/commander/domain/user"
	"github.com/vjftw/orchestrate/commander/persisters"
	"github.com/vjftw/orchestrate/commander/routers"
)

// OrchestrateApp - Orchestrate application struct
type OrchestrateApp struct {
	graph  *inject.Graph
	Router *routers.MuxRouter
}

type injectLogger struct {
}

func (l injectLogger) Debugf(format string, v ...interface{}) {
	// fmt.Println(fmt.Sprintf(format, v))
}

// NewOrchestrateApp - Initialise with Depencency Injection
func NewOrchestrateApp() *OrchestrateApp {
	orchestrateApp := OrchestrateApp{}
	orchestrateApp.graph = &inject.Graph{
		Logger: injectLogger{},
	}

	var userController user.Controller
	var authController auth.Controller
	var projectController project.Controller

	gormPersister := persisters.NewGORM(&user.User{}, &project.Project{})

	err := orchestrateApp.graph.Provide(
		&inject.Object{Name: "persister.gorm", Value: gormPersister},
		&inject.Object{Name: "user.manager", Value: user.NewManager()},
		&inject.Object{Name: "user.validator", Value: user.NewValidator()},
		&inject.Object{Name: "user.provider", Value: user.NewProvider()},
		&inject.Object{Name: "user.resolver", Value: user.NewResolver()},
		&inject.Object{Name: "auth.provider", Value: auth.NewProvider()},
		&inject.Object{Name: "project.manager", Value: project.NewManager()},
		&inject.Object{Name: "project.provider", Value: project.NewProvider()},
		&inject.Object{Name: "project.resolver", Value: project.NewResolver()},
		&inject.Object{Name: "project.validator", Value: project.NewValidator()},
		&inject.Object{
			Name:  "user.controller",
			Value: &userController,
		},
		&inject.Object{
			Name:  "auth.controller",
			Value: &authController,
		},
		&inject.Object{
			Name:  "project.controller",
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
