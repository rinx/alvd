package config

import (
	valdconfig "github.com/rinx/alvd/internal/config"
	"github.com/rinx/alvd/internal/log"
	ngt "github.com/rinx/alvd/pkg/vald/agent/ngt/config"
)

type Config struct {
	ServerAddress string

	AgentName string
	AgentPort uint

	NGTConfig *ngt.Data
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

func newNGTConfig() *ngt.Data {
	cfg := &ngt.Data{}

	cfg.Bind()

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	} else {
		cfg.Server = new(valdconfig.Servers)
	}

	if cfg.Observability != nil {
		cfg.Observability = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(valdconfig.Observability)
	}

	if cfg.NGT != nil {
		cfg.NGT = cfg.NGT.Bind()
	} else {
		cfg.NGT = new(valdconfig.NGT)
	}

	return cfg
}
