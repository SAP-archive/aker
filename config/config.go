// Package config defines the structure of Aker's configuration.
package config

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
