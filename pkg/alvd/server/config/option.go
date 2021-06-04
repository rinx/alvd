package config

import (
	"time"

	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/extension/lua"
)

type OptionFunc func(c *Config) error

func WithAgentEnabled(enabled bool) OptionFunc {
	return func(c *Config) error {
		c.AgentEnabled = enabled

		return nil
	}
}

func WithAgentOpts(opts *agent.Opts) OptionFunc {
	return func(c *Config) error {
		if opts != nil {
			c.AgentOpts = opts
		}

		return nil
	}
}

func WithAddrs(addrs []string) OptionFunc {
	return func(c *Config) error {
		c.Addrs = addrs

		return nil
	}
}

func WithGRPCHost(host string) OptionFunc {
	return func(c *Config) error {
		if host != "" {
			c.GRPCHost = host
		}

		return nil
	}
}

func WithGRPCPort(port uint) OptionFunc {
	return func(c *Config) error {
		c.GRPCPort = int(port)

		return nil
	}
}

func WithReplicas(n uint) OptionFunc {
	return func(c *Config) error {
		c.Replicas = int(n)

		return nil
	}
}

func WithCheckIndexInterval(s string) OptionFunc {
	return func(c *Config) error {
		dur, err := time.ParseDuration(s)
		if err != nil {
			return err
		}

		c.CheckIndexInterval = dur

		return nil
	}
}

func WithCreateIndexThreshold(n uint) OptionFunc {
	return func(c *Config) error {
		c.CreateIndexThreshold = int(n)

		return nil
	}
}

func WithSearchResultInterceptor(sri *lua.LFunction) OptionFunc {
	return func(c *Config) error {
		c.SearchResultInterceptor = sri

		return nil
	}
}

func WithSearchQueryInterceptor(sqi *lua.LFunction) OptionFunc {
	return func(c *Config) error {
		c.SearchQueryInterceptor = sqi

		return nil
	}
}

func WithInsertDataInterceptor(idi *lua.LFunction) OptionFunc {
	return func(c *Config) error {
		c.InsertDataInterceptor = idi

		return nil
	}
}
