package agent

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/log/level"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/runner"
	"github.com/rinx/alvd/pkg/alvd/extension/lua"
	"github.com/rinx/alvd/pkg/alvd/observability"

	cli "github.com/urfave/cli/v2"
)

type Opts struct {
	ServerAddresses        []string
	AgentName              string
	LogLevel               string
	Dimension              int
	DistanceType           string
	ObjectType             string
	CreationEdgeSize       int
	SearchEdgeSize         int
	BulkInsertChunkSize    int
	IndexPath              string
	IndexSelfcheckInterval string
	GRPCHost               string
	GRPCPort               uint
	MetricsHost            string
	MetricsPort            uint
	MetricsCollectInterval string
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "config",
		Value: "",
		Usage: "path to the config Lua file.",
	},
	&cli.StringFlag{
		Name:  "name",
		Value: "",
		Usage: "agent name (if not specified, uuid will be generated)",
	},
	&cli.StringSliceFlag{
		Name:  "server",
		Value: cli.NewStringSlice("0.0.0.0:8000"),
		Usage: "alvd server addresses",
	},
	&cli.StringFlag{
		Name:  "log-level",
		Value: "info",
		Usage: "log level",
	},
	&cli.IntFlag{
		Name:  "dimension",
		Value: 784,
		Usage: "dimension of vector",
	},
	&cli.StringFlag{
		Name:  "distance-type",
		Value: "l2",
		Usage: "distance type",
	},
	&cli.StringFlag{
		Name:  "object-type",
		Value: "float",
		Usage: "object type",
	},
	&cli.IntFlag{
		Name:  "creation-edge-size",
		Value: 10,
		Usage: "creation edge size",
	},
	&cli.IntFlag{
		Name:  "search-edge-size",
		Value: 40,
		Usage: "search edge size",
	},
	&cli.IntFlag{
		Name:  "bulk-insert-chunk-size",
		Value: 100,
		Usage: "bulk insert chunk size",
	},
	&cli.StringFlag{
		Name:  "index-path",
		Value: "",
		Usage: "index path (if not specified, in-memory mode will be enabled)",
	},
	&cli.StringFlag{
		Name:  "index-selfcheck-interval",
		Value: "24h",
		Usage: "selfcheck interval for alvd agent uncommitted index",
	},
	&cli.StringFlag{
		Name:  "grpc-host",
		Value: "0.0.0.0",
		Usage: "agent gRPC API host",
	},
	&cli.UintFlag{
		Name:  "grpc-port",
		Value: 8081,
		Usage: "agent gRPC API port",
	},
	&cli.StringFlag{
		Name:  "metrics-host",
		Value: "0.0.0.0",
		Usage: "metrics server host",
	},
	&cli.UintFlag{
		Name:  "metrics-port",
		Value: 9090,
		Usage: "metrics server port",
	},
	&cli.StringFlag{
		Name:  "metrics-collect-interval",
		Value: "5s",
		Usage: "interval for collecting metrics",
	},
}

func ParseConfig(c *cli.Context) (*Opts, error) {
	opts := ParseOpts(c)

	if config := c.String("config"); config != "" {
		err := lua.MapConfig(config, "agent", opts)
		if err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func ParseOpts(c *cli.Context) *Opts {
	return &Opts{
		AgentName:              c.String("name"),
		ServerAddresses:        c.StringSlice("server"),
		LogLevel:               c.String("log-level"),
		Dimension:              c.Int("dimension"),
		DistanceType:           c.String("distance-type"),
		ObjectType:             c.String("object-type"),
		CreationEdgeSize:       c.Int("creation-edge-size"),
		SearchEdgeSize:         c.Int("search-edge-size"),
		BulkInsertChunkSize:    c.Int("bulk-insert-chunk-size"),
		IndexPath:              c.String("index-path"),
		IndexSelfcheckInterval: c.String("index-selfcheck-interval"),
		GRPCHost:               c.String("grpc-host"),
		GRPCPort:               c.Uint("grpc-port"),
		MetricsHost:            c.String("metrics-host"),
		MetricsPort:            c.Uint("metrics-port"),
		MetricsCollectInterval: c.String("metrics-collect-interval"),
	}
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "Start agent",
		Flags: Flags,
		Action: func(c *cli.Context) error {
			opts, err := ParseConfig(c)
			if err != nil {
				return err
			}

			log.Init(log.WithLevel(level.Atol(opts.LogLevel).String()))
			log.Info("start alvd agent")

			cfg, err := ToConfig(opts)
			if err != nil {
				return err
			}

			ctx := context.Background()

			obs, err := observability.New(
				&observability.Config{
					MetricsHost:            opts.MetricsHost,
					MetricsPort:            opts.MetricsPort,
					MetricsCollectInterval: opts.MetricsCollectInterval,
				},
			)
			if err != nil {
				return err
			}

			err = obs.Start(ctx)
			if err != nil {
				return err
			}

			return Run(ctx, cfg)
		},
	}
}

func ToConfig(opts *Opts) (*config.Config, error) {
	cfg, err := config.New(
		config.WithAgentName(opts.AgentName),
		config.WithServerAddresses(opts.ServerAddresses),
		config.WithDimension(opts.Dimension),
		config.WithDistanceType(opts.DistanceType),
		config.WithObjectType(opts.ObjectType),
		config.WithCreationEdgeSize(opts.CreationEdgeSize),
		config.WithSearchEdgeSize(opts.SearchEdgeSize),
		config.WithBulkInsertChunkSize(opts.BulkInsertChunkSize),
		config.WithIndexPath(opts.IndexPath),
		config.WithAutoIndexCheckDuration(opts.IndexSelfcheckInterval),
		config.WithAutoIndexDurationLimit(opts.IndexSelfcheckInterval),
		config.WithAutoSaveIndexDuration(opts.IndexSelfcheckInterval),
		config.WithAutoIndexLength(1),
		config.WithProactiveGC(true),
		config.WithDefaultPoolSize(10000),
		config.WithDefaultRadius(-1.0),
		config.WithDefaultEpsilon(0.01),
		config.WithGRPCHost(opts.GRPCHost),
		config.WithGRPCPort(opts.GRPCPort),
	)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Run(ctx context.Context, cfg *config.Config) error {
	r, err := runner.New(cfg)
	if err != nil {
		return err
	}

	return r.Start(ctx)
}
