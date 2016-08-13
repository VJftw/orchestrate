package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/websocket"
	"github.com/vjftw/orchestrate/cadet/configuration"
)

func main() {
	// Read vars
	config := configuration.AutoConfiguration()

	// register cadet with commander
	url := fmt.Sprintf("http://%s/v1/cadets", config.CommanderAddress)
	type registerMessage struct {
		CadetGroupKey string `json:"cadetGroupKey"`
	}
	message := &registerMessage{
		CadetGroupKey: config.CadetGroupKey,
	}
	jsonBody, _ := json.Marshal(message)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(config)

	if err != nil {
		panic(err)
	}

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
	fmt.Println(resp)
	wsConn.WriteJSON(map[string]string{
		"hello": "world",
	})
}
