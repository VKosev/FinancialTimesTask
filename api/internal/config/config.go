package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Server struct {
		Port string `yaml:"port"`
	}
}

func New() (*Configuration, error) {
	config := &Configuration{}

	// Open config file
	file, err := os.Open("../../config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
