package server

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/log/level"
	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/runner"

	cli "github.com/urfave/cli/v2"
)

type Opts struct {
	AgentEnabled   bool
	ServerGRPCHost string
	ServerGRPCPort uint
	*agent.Opts
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
}

func ParseOpts(c *cli.Context) *Opts {
	return &Opts{
		AgentEnabled:   c.Bool("agent"),
		ServerGRPCHost: c.String("server-grpc-host"),
		ServerGRPCPort: c.Uint("server-grpc-port"),
		Opts:           agent.ParseOpts(c),
	}
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start server",
		Flags: append(Flags, agent.Flags...),
		Action: func(c *cli.Context) error {
			return Run(ParseOpts(c))
		},
	}
}

func Run(opts *Opts) error {
	log.Init(log.WithLevel(level.Atol(opts.LogLevel).String()))

	log.Info("start alvd server")

	cfg, err := config.New(
		config.WithAgentEnabled(opts.AgentEnabled),
		config.WithAddr(opts.ServerAddress),
		config.WithGRPCHost(opts.ServerGRPCHost),
		config.WithGRPCPort(opts.ServerGRPCPort),
	)
	if err != nil {
		return err
	}

	r, err := runner.New(cfg)
	if err != nil {
		return err
	}

	ctx := context.Background()

	return r.Start(ctx)
}
