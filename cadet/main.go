package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vjftw/orchestrate/cadet/configuration"
)

func main() {
	// Read vars
	config := configuration.AutoConfiguration()

	// register cadet with commander
	url := fmt.Sprintf("%s/v1/cadets", config.CommanderAddress)
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

	// connect to websocket with cadet key
}
