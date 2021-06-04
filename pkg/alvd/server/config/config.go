package config

import (
	"time"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/extension/lua"
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

	SearchResultInterceptor *lua.LFunction
	InsertDataInterceptor   *lua.LFunction
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
