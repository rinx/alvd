package runner

import (
	"context"

	"github.com/rinx/alvd/pkg/alvd/server/config"
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
	return nil
}
