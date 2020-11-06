package server

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/log/level"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/runner"

	cli "github.com/urfave/cli/v2"
)

type opts struct {
	logLevel string
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Usage: "log level",
			},
		},
		Action: func(c *cli.Context) error {
			return Run(&opts{
				logLevel: c.String("log-level"),
			})
		},
	}
}

func Run(opts *opts) error {
	log.Init(log.WithLevel(level.Atol(opts.logLevel).String()))

	log.Info("start alvd server")

	cfg, err := config.New()
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
