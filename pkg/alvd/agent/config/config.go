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
	cfg := &valdconfig.NGT{
		IndexPath:               "",
		Dimension:               784,
		BulkInsertChunkSize:     10,
		DistanceType:            "l2",
		ObjectType:              "float",
		CreationEdgeSize:        20,
		SearchEdgeSize:          10,
		AutoIndexDurationLimit:  "24h",
		AutoIndexCheckDuration:  "30m",
		AutoSaveIndexDuration:   "31m",
		AutoIndexLength:         100,
		InitialDelayMaxDuration: "30s",
		EnableInMemoryMode:      true,
		DefaultPoolSize:         10000,
		DefaultRadius:           -1.0,
		DefaultEpsilon:          0.1,
		MinLoadIndexTimeout:     "3m",
		MaxLoadIndexTimeout:     "10m",
		LoadIndexTimeoutFactor:  "1ms",
		EnableProactiveGC:       true,
		VQueue: &valdconfig.VQueue{
			InsertBufferSize:     100,
			InsertBufferPoolSize: 100,
			DeleteBufferSize:     10000,
			DeleteBufferPoolSize: 5000,
		},
	}

	cfg.Bind()

	return cfg
}
