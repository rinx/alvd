package config

import (
	"github.com/kpango/fuid"
	valdconfig "github.com/rinx/alvd/internal/config"
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
			c.NGTConfig.NGT.Dimension = dimension
		}

		return nil
	}
}

func WithDistanceType(dt string) OptionFunc {
	return func(c *Config) error {
		if dt != "" {
			c.NGTConfig.NGT.DistanceType = dt
		}

		return nil
	}
}

func WithObjectType(ot string) OptionFunc {
	return func(c *Config) error {
		if ot != "" {
			c.NGTConfig.NGT.ObjectType = ot
		}

		return nil
	}
}

func WithCreationEdgeSize(size int) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.NGT.CreationEdgeSize = size
		}

		return nil
	}
}

func WithSearchEdgeSize(size int) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.NGT.SearchEdgeSize = size
		}

		return nil
	}
}

func WithBulkInsertChunkSize(size int) OptionFunc {
	return func(c *Config) error {
		if size != 0 {
			c.NGTConfig.NGT.BulkInsertChunkSize = 0
		}

		return nil
	}
}

func WithIndexPath(path string) OptionFunc {
	return func(c *Config) error {
		if path == "" {
			c.NGTConfig.NGT.EnableInMemoryMode = true
			return nil
		}

		c.NGTConfig.NGT.EnableInMemoryMode = false
		c.NGTConfig.NGT.IndexPath = path

		return nil
	}
}

func WithGRPCServer(enable bool, host string, port uint) OptionFunc {
	return func(c *Config) error {
		if enable {
			c.NGTConfig.Server.Servers = append(
				c.NGTConfig.Server.Servers,
				&valdconfig.Server{
					Name: "grpc",
					Host: host,
					Port: port,
					Mode: "GRPC",
				},
			)
			c.AgentPort = port
		}

		return nil
	}
}
