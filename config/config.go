package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Server   ServerConfig    `json:"server"`
	Handlers []HandlerConfig `json:"handlers"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HandlerConfig struct {
	Path    string         `json:"path"`
	Filters []FilterConfig `json:"filters"`
}

type FilterConfig struct {
	PluginName   string       `json:"plugin_name"`
	PluginConfig PluginConfig `json:"plugin_config"`
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
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
