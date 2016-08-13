package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/gorilla/websocket"
	"github.com/vjftw/orchestrate/cadet/configuration"
	"github.com/vjftw/orchestrate/cadet/connection"
	"github.com/vjftw/orchestrate/cadet/node"
	"github.com/vjftw/orchestrate/cadet/registration"
)

func main() {
	// Read vars
	config := configuration.AutoConfiguration()

	// Register Cadet using vars
	registration.Register(config)

	// Print Docker containers
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Println(c.Names)
	}

	// connect to websocket with cadet key
	dialer := websocket.Dialer{}
	wsurl := fmt.Sprintf("ws://%s/v1/cadets/%s/ws", config.CommanderAddress, config.CadetUUID)
	ws, _, err := dialer.Dial(wsurl, nil)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("[websocket] Connected to %s", wsurl))

	log.Println(fmt.Sprintf("[websocket] Authenticating with Cadet Key: %s", config.CadetKey))
	n := node.NewNode()
	ws.WriteJSON(connection.NewMessage(config.CadetKey, n))

	wsJSON := node.NewNode()
	err = ws.ReadJSON(wsJSON)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Closing WebSocket")
		ws.Close()
		return
	}
	log.Println(fmt.Sprintf("[websocket] Authenticated."))

	monitorConn(ws)
}

func monitorConn(ws *websocket.Conn) {
	for {
		wsJSON := node.NewNode()
		err := ws.ReadJSON(wsJSON)
		if err != nil {
			log.Println(fmt.Sprintf("[websocket] Closing Websocket - %s", err))
			ws.Close()
			return
		}

		fmt.Println(wsJSON)
	}
}
