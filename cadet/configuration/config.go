package configuration

import (
	"fmt"
	"log"
	"os"
)

// Configuration - Holds configuration for a Cadet
type Configuration struct {
	CommanderAddress string `json:"commanderAddress"`
	CadetGroupKey    string `json:"cadetGroupKey"`
	CadetUUID        string `json:"uuid"`
	CadetKey         string `json:"key"`
}

// AutoConfiguration - Automatically finds configuration based on environment
func AutoConfiguration() *Configuration {
	host := os.Getenv("CADET_HOST")
	log.Println(fmt.Sprintf("[configuration] Running on: %s", host))

	switch host {
	case "AWS":
		return NewAWS()
	case "GCP":
		return NewGCP()
	case "ENV":
		return NewEnv()
	}

	return nil
}
