package agent

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/log/level"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/runner"

	cli "github.com/urfave/cli/v2"
)

type Opts struct {
	ServerAddress string
	AgentName     string
	LogLevel      string
	Dimension     int
	DistanceType  string
	ObjectType    string
	RESTEnabled   bool
	RESTHost      string
	RESTPort      uint
	GRPCEnabled   bool
	GRPCHost      string
	GRPCPort      uint
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "name",
		Value: "",
		Usage: "agent name (if not specified, uuid will be generated)",
	},
	&cli.StringFlag{
		Name:  "server",
		Value: "0.0.0.0:8000",
		Usage: "server address",
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
	&cli.BoolFlag{
		Name:  "rest",
		Value: false,
		Usage: "rest server enabled",
	},
	&cli.StringFlag{
		Name:  "rest-host",
		Value: "0.0.0.0",
		Usage: "rest server host",
	},
	&cli.UintFlag{
		Name:  "rest-port",
		Value: 8080,
		Usage: "rest server port",
	},
	&cli.BoolFlag{
		Name:  "grpc",
		Value: true,
		Usage: "grpc server enabled",
	},
	&cli.StringFlag{
		Name:  "grpc-host",
		Value: "0.0.0.0",
		Usage: "grpc server host",
	},
	&cli.UintFlag{
		Name:  "grpc-port",
		Value: 8081,
		Usage: "grpc server port",
	},
}

func ParseOpts(c *cli.Context) *Opts {
	return &Opts{
		AgentName:     c.String("name"),
		ServerAddress: c.String("server"),
		LogLevel:      c.String("log-level"),
		Dimension:     c.Int("dimension"),
		DistanceType:  c.String("distance-type"),
		ObjectType:    c.String("object-type"),
		RESTEnabled:   c.Bool("rest"),
		RESTHost:      c.String("rest-host"),
		RESTPort:      c.Uint("rest-port"),
		GRPCEnabled:   c.Bool("grpc"),
		GRPCHost:      c.String("grpc-host"),
		GRPCPort:      c.Uint("grpc-port"),
	}
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "Start agent",
		Flags: Flags,
		Action: func(c *cli.Context) error {
			return Run(ParseOpts(c))
		},
	}
}

func Run(opts *Opts) error {
	log.Init(log.WithLevel(level.Atol(opts.LogLevel).String()))

	log.Info("start alvd agent")

	cfg, err := config.New(
		config.WithAgentName(opts.AgentName),
		config.WithServerAddress(opts.ServerAddress),
		config.WithDimension(opts.Dimension),
		config.WithDistanceType(opts.DistanceType),
		config.WithObjectType(opts.ObjectType),
		config.WithRESTServer(opts.RESTEnabled, opts.RESTHost, opts.RESTPort),
		config.WithGRPCServer(opts.GRPCEnabled, opts.GRPCHost, opts.GRPCPort),
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
