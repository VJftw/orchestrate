package configuration

import "os"

// NewEnv - Returns a configuration from system environment variables
func NewEnv() *Configuration {
	return &Configuration{
		CommanderAddress: os.Getenv("CADET_COMMANDER_ADDRESS"),
		CadetGroupKey:    os.Getenv("CADET_GROUP_KEY"),
	}
}
