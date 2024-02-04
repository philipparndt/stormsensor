package config

import (
	"encoding/json"
	"github.com/philipparndt/go-logger"
	"github.com/philipparndt/mqtt-gateway/config"
	"os"
)

type Port struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type Storm struct {
 	WindSpeed        float64    `json:"windSpeed"`
	ResetTimeSeconds int    `json:"resetTimeSeconds"`
	WindTopic        string `json:"windTopic"`
}

type Config struct {
	MQTT     config.MQTTConfig `json:"mqtt"`
	Storm    Storm             `json:"storm"`
	LogLevel string            `json:"loglevel,omitempty"`
}

func LoadConfig(file string) (Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		logger.Error("Error reading config file", err)
		return Config{}, err
	}

	data = config.ReplaceEnvVariables(data)

	// Create a Config object
	var cfg Config

	// Unmarshal the JSON data into the Config object
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("Unmarshalling JSON:", err)
		return Config{}, err
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}
