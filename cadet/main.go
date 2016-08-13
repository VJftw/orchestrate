package main

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
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

	fmt.Println(config)

	// Print Docker containers
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	containers, _ := client.ListContainers(docker.ListContainersOptions{})
	for _, container := range containers {
		fmt.Println(container.Command, container.Image, container.Names, container.Ports, container.Status)
	}

	// connect to websocket with cadet key
	dialer := websocket.Dialer{}
	wsurl := fmt.Sprintf("ws://%s/v1/cadets/%s/ws", config.CommanderAddress, config.CadetUUID)
	wsConn, resp, err := dialer.Dial(wsurl, nil)
	if err != nil {
		panic(err)
	}
	n := node.NewNode()
	fmt.Println(resp)

	wsConn.WriteJSON(connection.NewMessage(config.CadetKey, n))

	monitorConn(wsConn)
}

func monitorConn(ws *websocket.Conn) {
	for {
		wsJSON := node.NewNode()
		err := ws.ReadJSON(wsJSON)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing WebSocket")
			ws.Close()
			return
		}

		fmt.Println(wsJSON)
	}
}
