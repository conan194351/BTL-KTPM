package config

import (
	"github.com/pkg/errors"
	"time"
)

type Server struct {
	Host        string `env:"HOST" envDefault:"localhost"`
	Port        string `env:"PORT" envDefault:"8080"`
	SecretKey   string `env:"SECRET_KEY" envDefault:"serversecretkey"`
	ReadTimeout string `env:"READ_TIMEOUT" envDefault:"10s"`
}

func (s Server) GetAddr() string {
	return s.Host + ":" + s.Port
}

func (s Server) GetReadTimeout() time.Duration {
	duration, err := time.ParseDuration(s.ReadTimeout)
	if err != nil {
		panic(errors.Wrap(err, "Failed to parse read timeout"))
	}
	return duration
}
