package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL"  required:"true" default:"INFO"`
	Folder string `envconfig:"FOLDER" required:"true" default:"logs"`
}

func NewConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return nil, fmt.Errorf("process envconfig: %w", err)
	}

	return &config, nil
}

// "...Must" used because there is no sense in handling the error in app initialization phase
// It's more convenient to panic and fix environmental variables and then restart the app

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}
	return config
}
