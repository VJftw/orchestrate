package configuration

import "os"

// Holds configuration for a Cadet
type Configuration struct {
	CommanderAddress string `json:"commanderAddress"`
	CadetGroupKey    string `json:"cadetGroupKey"`
	CadetUUID        string `json:"uuid"`
	CadetKey         string `json:"key"`
}

func AutoConfiguration() *Configuration {
	switch os.Getenv("CADET_HOST") {
	case "AWS":
		return NewAWS()
	case "GCP":
		return NewGCP()
	case "ENV":
		return NewEnv()
	}

	return nil
}
