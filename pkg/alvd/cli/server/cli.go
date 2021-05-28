package server

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/log/level"
	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/extension/lua"
	"github.com/rinx/alvd/pkg/alvd/observability"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/runner"

	cli "github.com/urfave/cli/v2"
)

type Opts struct {
	AgentEnabled bool

	LogLevel string

	ServerAddresses []string
	ServerGRPCHost  string
	ServerGRPCPort  uint

	MetricsHost            string
	MetricsPort            uint
	MetricsCollectInterval string

	Replicas             uint
	CheckIndexInterval   string
	CreateIndexThreshold uint

	EgressFilter *lua.LFunction
}

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "agent",
		Value: true,
		Usage: "agent enabled",
	},
	&cli.StringFlag{
		Name:  "server-grpc-host",
		Value: "0.0.0.0",
		Usage: "alvd server gRPC API host",
	},
	&cli.UintFlag{
		Name:  "server-grpc-port",
		Value: 8080,
		Usage: "alvd server gRPC API host",
	},
	&cli.UintFlag{
		Name:  "replicas",
		Value: 3,
		Usage: "number of index replicas",
	},
	&cli.StringFlag{
		Name:  "check-index-interval",
		Value: "5s",
		Usage: "check interval for alvd agent index",
	},
	&cli.UintFlag{
		Name:  "create-index-threshold",
		Value: 100,
		Usage: "number of data to trigger create index",
	},
}

func ParseConfig(c *cli.Context) (*Opts, error) {
	opts := ParseOpts(c)

	if config := c.String("config"); config != "" {
		err := lua.MapConfig(config, "server", opts)
		if err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func ParseOpts(c *cli.Context) *Opts {
	return &Opts{
		AgentEnabled:           c.Bool("agent"),
		LogLevel:               c.String("log-level"),
		ServerAddresses:        c.StringSlice("server"),
		ServerGRPCHost:         c.String("server-grpc-host"),
		ServerGRPCPort:         c.Uint("server-grpc-port"),
		Replicas:               c.Uint("replicas"),
		CheckIndexInterval:     c.String("check-index-interval"),
		CreateIndexThreshold:   c.Uint("create-index-threshold"),
		MetricsHost:            c.String("metrics-host"),
		MetricsPort:            c.Uint("metrics-port"),
		MetricsCollectInterval: c.String("metrics-collect-interval"),
	}
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start server",
		Flags: append(Flags, agent.Flags...),
		Action: func(c *cli.Context) error {
			opts, err := ParseConfig(c)
			if err != nil {
				return err
			}

			agentOpts, err := agent.ParseConfig(c)
			if err != nil {
				return err
			}

			log.Init(log.WithLevel(level.Atol(opts.LogLevel).String()))
			log.Info("start alvd server")

			cfg, err := ToConfig(opts, agentOpts)
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

func ToConfig(opts *Opts, agentOpts *agent.Opts) (*config.Config, error) {
	cfg, err := config.New(
		config.WithAgentEnabled(opts.AgentEnabled),
		config.WithAgentOpts(agentOpts),
		config.WithAddrs(opts.ServerAddresses),
		config.WithGRPCHost(opts.ServerGRPCHost),
		config.WithGRPCPort(opts.ServerGRPCPort),
		config.WithReplicas(opts.Replicas),
		config.WithCheckIndexInterval(opts.CheckIndexInterval),
		config.WithCreateIndexThreshold(opts.CreateIndexThreshold),
		config.WithEgressFilter(opts.EgressFilter),
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
