package cadet

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"github.com/vjftw/orchestrate/commander/domain/cadetGroup"
)

// Controller - Handles actions that can be performed on Cadets
type Controller struct {
	render            *render.Render
	CadetResolver     Resolver           `inject:"cadet.resolver"`
	CadetGroupManager cadetGroup.Manager `inject:"cadetGroup.manager"`
	CadetManager      Manager            `inject:"cadet.manager"`
	WSLog             *log.Logger        `inject:"logger.ws"`
}

// Setup - Sets up the controller on the router and a renderer
func (c Controller) Setup(router *mux.Router, renderer *render.Render) {
	c.render = renderer

	router.
		HandleFunc("/v1/cadets", c.postHandler).
		Methods("POST")

	router.
		HandleFunc("/v1/cadets/{cadetUUID}/ws", c.wsHandler).
		Methods("GET")
}

func (c Controller) postHandler(w http.ResponseWriter, r *http.Request) {
	// use resolver to get cadetgroupkey from POST body
	cadetGroupKey, err := c.CadetResolver.KeyFromRequest(r.Body)
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// find CadetGroup by key
	cadetGroup, err := c.CadetGroupManager.FindByKey(cadetGroupKey)
	if err != nil {
		c.render.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	// create new Cadet
	cadet := c.CadetManager.NewForCadetGroup(cadetGroup)

	// add new key to Cadet
	cadet.UUID = uuid.NewV4().String()
	secureRandom := make([]byte, 10)
	rand.Read(secureRandom)
	keyBytes := sha512.Sum512_256(secureRandom)
	cadet.Key = hex.EncodeToString(keyBytes[:sha512.Size256])

	// persist
	c.CadetManager.Save(cadet)

	// return Cadet
	c.render.JSON(w, http.StatusCreated, cadet)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c Controller) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Find Cadet
	cadetUUID := mux.Vars(r)["cadetUUID"]
	cadet, err := c.CadetManager.FindByUUID(cadetUUID)

	if err != nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go c.monitorCadet(ws, cadet)
}

func (c Controller) monitorCadet(ws *websocket.Conn, cadet *Cadet) {
	for {
		message := NewMessage()
		err := ws.ReadJSON(message)
		if err != nil {
			c.WSLog.Println(fmt.Sprintf("Closing Websocket: %s", err))
			ws.Close()
			return
		}

		c.WSLog.Println(fmt.Sprintf("Recieved: %v", message))

		if cadet.Key != message.Key {
			ws.WriteJSON(map[string]string{
				"error": "Invalid Key",
			})
			c.WSLog.Println("Closing Websocket: Invalid Key")
			ws.Close()
			return
		}
		// Add containers to message object

		// send Message
		ws.WriteJSON(NewMessage())
	}
}
