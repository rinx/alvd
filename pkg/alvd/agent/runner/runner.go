package runner

import (
	"context"

	valdrunner "github.com/rinx/alvd/internal/runner"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/daemon"
	"github.com/rinx/alvd/pkg/vald/agent/ngt/usecase"
)

type runner struct {
	cfg *config.Config

	daemon daemon.Daemon
}

type Runner interface {
	Start(ctx context.Context) error
}

func New(cfg *config.Config) (Runner, error) {
	d, err := daemon.New(cfg)
	if err != nil {
		return nil, err
	}

	return &runner{
		cfg:    cfg,
		daemon: d,
	}, nil
}

func (r *runner) Start(ctx context.Context) error {
	if r.cfg.NGTConfig == nil {
		return nil
	}

	vr, err := usecase.New(r.cfg.NGTConfig)
	if err != nil {
		return err
	}

	err = r.daemon.Start(ctx)
	if err != nil {
		return err
	}
	defer r.daemon.Close()

	return valdrunner.Run(ctx, vr, "alvd-agent")
}
