package configuration

import (
	"fmt"
	"log"
	"os"
)

// NewEnv - Returns a configuration from system environment variables
func NewEnv() *Configuration {
	commanderAddress := os.Getenv("CADET_COMMANDER_ADDRESS")
	cadetGroupKey := os.Getenv("CADET_GROUP_KEY")
	log.Println(fmt.Sprintf("[configuration] Commander Address: %s", commanderAddress))
	log.Println(fmt.Sprintf("[configuration] Cadet Group Key: %s", cadetGroupKey))
	return &Configuration{
		CommanderAddress: commanderAddress,
		CadetGroupKey:    cadetGroupKey,
	}
}
