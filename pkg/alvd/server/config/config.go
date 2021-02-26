package config

import (
	"time"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/cli/agent"
)

type Config struct {
	AgentEnabled bool
	AgentOpts    *agent.Opts

	Addrs    []string
	GRPCHost string
	GRPCPort int

	Replicas             int
	CheckIndexInterval   time.Duration
	CreateIndexThreshold int
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
