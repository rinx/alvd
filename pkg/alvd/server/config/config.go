package config

import (
	"github.com/rinx/alvd/internal/log"
)

type Config struct {
	AgentEnabled bool

	Addr string
}

func New(opts ...OptionFunc) (*Config, error) {
	cfg := &Config{}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			log.Errorf("%s", err)

			return nil, err
		}
	}

	return cfg, nil
}
