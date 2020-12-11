package config

import (
	"github.com/kpango/fuid"
)

type OptionFunc func(c *Config) error

func WithAgentName(name string) OptionFunc {
	return func(c *Config) error {
		if name != "" {
			c.AgentName = name
		} else {
			c.AgentName = fuid.String()
		}

		return nil
	}
}

func WithServerAddress(addr string) OptionFunc {
	return func(c *Config) error {
		if addr != "" {
			c.ServerAddress = addr
		}

		return nil
	}
}

func WithDimension(dimension int) OptionFunc {
	return func(c *Config) error {
		if dimension > 2 {
			c.NGTConfig.Dimension = dimension
		}

		return nil
	}
}

func WithDistanceType(dt string) OptionFunc {
	return func(c *Config) error {
		if dt != "" {
			c.NGTConfig.DistanceType = dt
		}

		return nil
	}
}

func WithObjectType(ot string) OptionFunc {
	return func(c *Config) error {
		if ot != "" {
			c.NGTConfig.ObjectType = ot
		}

		return nil
	}
}

func WithCreationEdgeSize(size int) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.CreationEdgeSize = size
		}

		return nil
	}
}

func WithSearchEdgeSize(size int) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.SearchEdgeSize = size
		}

		return nil
	}
}

func WithBulkInsertChunkSize(size int) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.BulkInsertChunkSize = 0
		}

		return nil
	}
}

func WithIndexPath(path string) OptionFunc {
	return func(c *Config) error {
		if path == "" {
			c.NGTConfig.EnableInMemoryMode = true
			return nil
		}

		c.NGTConfig.EnableInMemoryMode = false
		c.NGTConfig.IndexPath = path

		return nil
	}
}

func WithAutoIndexCheckDuration(s string) OptionFunc {
	return func(c *Config) error {
		c.NGTConfig.AutoIndexCheckDuration = s

		return nil
	}
}

func WithAutoIndexDurationLimit(s string) OptionFunc {
	return func(c *Config) error {
		c.NGTConfig.AutoIndexDurationLimit = s

		return nil
	}
}

func WithAutoSaveIndexDuration(s string) OptionFunc {
	return func(c *Config) error {
		c.NGTConfig.AutoSaveIndexDuration = s

		return nil
	}
}

func WithAutoIndexLength(n int) OptionFunc {
	return func(c *Config) error {
		if n > 0 {
			c.NGTConfig.AutoIndexLength = n
		}

		return nil
	}
}

func WithProactiveGC(enabled bool) OptionFunc {
	return func(c *Config) error {
		c.NGTConfig.EnableProactiveGC = enabled

		return nil
	}
}

func WithDefaultPoolSize(size uint32) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.DefaultPoolSize = size
		}

		return nil
	}
}

func WithDefaultRadius(r float32) OptionFunc {
	return func(c *Config) error {
		c.NGTConfig.DefaultRadius = r

		return nil
	}
}

func WithDefaultEpsilon(e float32) OptionFunc {
	return func(c *Config) error {
		c.NGTConfig.DefaultEpsilon = e

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
		if port != 0 {
			c.GRPCPort = int(port)
		}

		return nil
	}
}
