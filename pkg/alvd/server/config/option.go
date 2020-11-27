package config

import "time"

type OptionFunc func(c *Config) error

func WithAgentEnabled(enabled bool) OptionFunc {
	return func(c *Config) error {
		c.AgentEnabled = enabled

		return nil
	}
}

func WithAddr(addr string) OptionFunc {
	return func(c *Config) error {
		if addr != "" {
			c.Addr = addr
		}

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
