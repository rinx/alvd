package config

import (
	valdconfig "github.com/rinx/alvd/internal/config"
)

type OptionFunc func(c *Config) error

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

func WithRESTServer(enable bool, port uint) OptionFunc {
	return func(c *Config) error {
		if enable {
			c.NGTConfig.Server.Servers = append(
				c.NGTConfig.Server.Servers,
				&valdconfig.Server{
					Name: "rest",
					Host: "0.0.0.0",
					Port: port,
					Mode: "REST",
				},
			)
		}

		return nil
	}
}

func WithGRPCServer(enable bool, port uint) OptionFunc {
	return func(c *Config) error {
		if enable {
			c.NGTConfig.Server.Servers = append(
				c.NGTConfig.Server.Servers,
				&valdconfig.Server{
					Name: "grpc",
					Host: "0.0.0.0",
					Port: port,
					Mode: "GRPC",
				},
			)
		}

		return nil
	}
}
