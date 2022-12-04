package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Configuration represents the configured properties as a struct
type Configuration struct {
	Server struct {
		Port string `yaml:"port"`
	}
}

// New reads the config.yaml file and returns a pointer to new instance of Configuration.
// Returns an error if it fails to read from the config.yaml file.
func New() (*Configuration, error) {
	config := &Configuration{}

	file, err := os.Open("../../config.yaml")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	d := yaml.NewDecoder(file)

	err = d.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
