package config

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	ServerHost string `default:"0.0.0.0"`
	ServerPort string `default:"8080"`
}

func LoadConfig() (*AppConfig, error) {
	appConfig := &AppConfig{}
	err := envconfig.Process("et", appConfig)
	if err != nil {
		return nil, errors.New("unable to load config")
	}

	return appConfig, nil
}
