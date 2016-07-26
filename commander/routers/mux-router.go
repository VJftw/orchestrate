package routers

// MuxRouter - The application router
import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type MuxRouter struct {
	Router  *mux.Router
	Render  *render.Render
	Handler http.Handler
}

type Routable interface {
	Setup(*mux.Router, *render.Render)
}

func NewMuxRouter(controllers []Routable, logging bool) *MuxRouter {
	muxRouter := MuxRouter{}

	muxRouter.Render = render.New()
	muxRouter.Router = mux.NewRouter()

	n := negroni.New()
	if logging {
		n.Use(negroni.NewLogger())
	}
	n.Use(negroni.NewRecovery())

	n.UseHandler(muxRouter.Router)

	muxRouter.Handler = n

	for _, controller := range controllers {
		controller.Setup(muxRouter.Router, muxRouter.Render)
	}

	return &muxRouter
}
