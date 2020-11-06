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
	LogLevel     string
	Dimension    int
	DistanceType string
	ObjectType   string
	RESTEnabled  bool
	RESTPort     uint
	GRPCEnabled  bool
	GRPCPort     uint
}

var Flags = []cli.Flag{
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
	&cli.UintFlag{
		Name:  "grpc-port",
		Value: 8081,
		Usage: "grpc server port",
	},
}

func ParseOpts(c *cli.Context) *Opts {
	return &Opts{
		LogLevel:     c.String("log-level"),
		Dimension:    c.Int("dimension"),
		DistanceType: c.String("distance-type"),
		ObjectType:   c.String("object-type"),
		RESTEnabled:  c.Bool("rest"),
		RESTPort:     c.Uint("rest-port"),
		GRPCEnabled:  c.Bool("grpc"),
		GRPCPort:     c.Uint("grpc-port"),
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
		config.WithDimension(opts.Dimension),
		config.WithDistanceType(opts.DistanceType),
		config.WithObjectType(opts.ObjectType),
		config.WithRESTServer(opts.RESTEnabled, opts.RESTPort),
		config.WithGRPCServer(opts.GRPCEnabled, opts.GRPCPort),
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
