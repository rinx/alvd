package runner

import (
	"context"

	valdrunner "github.com/rinx/alvd/internal/runner"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/tunnel"
	"github.com/rinx/alvd/pkg/vald/agent/ngt/usecase"
)

type runner struct {
	cfg *config.Config
}

type Runner interface {
	Start(ctx context.Context) error
}

func New(cfg *config.Config) (Runner, error) {
	return &runner{
		cfg: cfg,
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

	tun := tunnel.Connect(ctx, r.cfg.ServerAddress)
	defer tun.Close()

	return valdrunner.Run(ctx, vr, "alvd-agent")
}
