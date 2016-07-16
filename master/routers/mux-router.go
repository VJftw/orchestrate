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

func NewMuxRouter() *MuxRouter {
	muxRouter := MuxRouter{}

	muxRouter.Router = mux.NewRouter()

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())

	n.UseHandler(muxRouter.Router)

	muxRouter.Handler = n

	muxRouter.Render = render.New()

	return &muxRouter
}
