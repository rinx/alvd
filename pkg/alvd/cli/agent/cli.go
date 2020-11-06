package agent

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/log/level"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/runner"

	cli "github.com/urfave/cli/v2"
)

type opts struct {
	logLevel     string
	dimension    int
	distanceType string
	objectType   string
	restEnabled  bool
	restPort     uint
	grpcEnabled  bool
	grpcPort     uint
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "Start agent",
		Flags: []cli.Flag{
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
		},
		Action: func(c *cli.Context) error {
			return Run(&opts{
				logLevel:     c.String("log-level"),
				dimension:    c.Int("dimension"),
				distanceType: c.String("distance-type"),
				objectType:   c.String("object-type"),
				restEnabled:  c.Bool("rest"),
				restPort:     c.Uint("rest-port"),
				grpcEnabled:  c.Bool("grpc"),
				grpcPort:     c.Uint("grpc-port"),
			})
		},
	}
}

func Run(opts *opts) error {
	log.Init(log.WithLevel(level.Atol(opts.logLevel).String()))

	log.Info("start alvd agent")

	cfg, err := config.New(
		config.WithDimension(opts.dimension),
		config.WithDistanceType(opts.distanceType),
		config.WithObjectType(opts.objectType),
		config.WithRESTServer(opts.restEnabled, opts.restPort),
		config.WithGRPCServer(opts.grpcEnabled, opts.grpcPort),
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
