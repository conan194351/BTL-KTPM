package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/pkg/errors"
)

type Config struct {
	App      App      `envPrefix:"APP_"`
	Server   Server   `envPrefix:"SERVER_"`
	Database Database `envPrefix:"DATABASE_"`
}

var config Config

func init() {
	err := env.Parse(&config)
	if err != nil {
		panic(errors.Wrap(err, "Failed to init config"))
	}

	timezone := config.App.Timezone

	if timezone == "" {
		timezone = "UTC"
	}

	if err := SetTimezone(timezone); err != nil {
		panic(errors.Wrap(err, "Failed to set timezone"))
	}
}

func GetConfig() Config {
	return config
}

func SetTimezone(timezone string) error {
	return os.Setenv("TZ", timezone)
}
