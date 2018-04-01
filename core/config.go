package core

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// Config configures the client.
type Config struct {
	Token         string    `yaml:"token"`
	DefaultPrefix string    `yaml:"prefix"`
	DatabasePath  string    `yaml:"database"`
	Shards        int       `yaml:"shards"`
	Sentry        bool      `yaml:"-"`
	Keys          KeyConfig `yaml:"keys"`
}

// KeyConfig contains the API keys.
type KeyConfig struct {
	Sentry         string `yaml:"sentry"`
	OpenWeatherMap string `yaml:"openweathermap"`
	Google         string `yaml:"google"`
	ChatEngine     string `yaml:"chatengine"`
}

// LoadConfig loads a client config from the given path.
func LoadConfig(path string) (Config, error) {
	config := Config{
		Shards:        1,
		DefaultPrefix: "!",
		DatabasePath:  "data/db",
		Sentry:        false,
		Keys:          KeyConfig{},
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
