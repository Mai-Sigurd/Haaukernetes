package utils

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Settings struct {
	Endpoint string `yaml:"endpoint"`
	Subnet   string `yaml:"subnet"`
}

func ReadYaml(filename string) Settings {
	// Load the file
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var s Settings

	// Unmarshal input YAML file into empty Settings
	if err := yaml.Unmarshal(f, &s); err != nil {
		ErrorLogger.Fatal(err)
	}

	// : needed for the wireguard config
	s.Endpoint = s.Endpoint + ":"

	return s
}
