package config

import (
	valdconfig "github.com/rinx/alvd/internal/config"
	"github.com/rinx/alvd/internal/log"
)

type Config struct {
	ServerAddresses []string

	AgentName string

	GRPCHost string
	GRPCPort int

	NGTConfig *valdconfig.NGT
}

func New(opts ...OptionFunc) (*Config, error) {
	cfg := &Config{
		NGTConfig: newNGTConfig(),
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			log.Errorf("%s", err)

			return nil, err
		}
	}

	return cfg, nil
}

func newNGTConfig() *valdconfig.NGT {
	cfg := &valdconfig.NGT{}

	cfg.Bind()

	return cfg
}
