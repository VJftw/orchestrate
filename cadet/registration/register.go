package registration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vjftw/orchestrate/cadet/configuration"
)

func Register(c *configuration.Configuration) {
	// register cadet with commander
	url := fmt.Sprintf("http://%s/v1/cadets", c.CommanderAddress)
	type registerMessage struct {
		CadetGroupKey string `json:"cadetGroupKey"`
	}
	message := &registerMessage{
		CadetGroupKey: c.CadetGroupKey,
	}
	jsonBody, _ := json.Marshal(message)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(c)
	if err != nil {
		panic(err)
	}
}
