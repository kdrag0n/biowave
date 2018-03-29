package core

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
)

// Config configures the client.
type Config struct {
	Token         string    `json:"token"`
	DefaultPrefix string    `json:"prefix"`
	Shards        int       `json:"shards"`
	Sentry        bool      `json:"-"`
	Keys          KeyConfig `json:"keys"`
}

// KeyConfig contains the API keys.
type KeyConfig struct {
	Sentry         string `json:"sentry"`
	OpenWeatherMap string `json:"openweathermap"`
	Google         string `json:"google"`
	ChatEngine     string `json:"chatengine"`
}

// LoadConfig loads a client config from the given path.
func LoadConfig(path string) (Config, error) {
	config := Config{
		Shards:        1,
		DefaultPrefix: "!",
		Sentry:        false,
		Keys:          KeyConfig{},
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("Error loading config", zap.Error(err))
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		logger.Error("Error loading config", zap.Error(err))
		return config, err
	}

	return config, nil
}
