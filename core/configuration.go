package core

import (
	"encoding/json"
	"log"
	"os"
)

type StateControlConfiguration struct {
	Instances   *map[string]string
	Autoscaling *map[string]string
}

type Configuration struct {
	Terminate *StateControlConfiguration
	Stop      *StateControlConfiguration
}

// ReadConfigiguration - gets configuration file from the specified location and
// applies the default values if needed
func ReadConfigiguration(filename string) Configuration {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Printf("Can't read configuration file %s : %s", filename, err)
	}
	return configuration
}
