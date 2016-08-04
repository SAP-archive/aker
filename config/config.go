// Package config defines the structure of Aker's configuration.
package config

import (
	"io"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

type Config struct {
	Server    ServerConfig `yaml:"server"`
	Endpoints []Endpoint   `yaml:"endpoints"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Endpoint struct {
	Path    string            `yaml:"path"`
	Audit   bool              `yaml:"audit"`
	Plugins []PluginReference `yaml:"plugins"`
}

type PluginReference struct {
	Name   string       `yaml:"name"`
	Config PluginConfig `yaml:"configuration"`
}

type PluginConfig map[string]interface{}

func LoadFromFile(name string) (Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	return loadFromReader(file)
}

func loadFromReader(reader io.Reader) (Config, error) {
	config := Config{}
	decoder := candiedyaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
