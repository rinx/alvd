package config

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
