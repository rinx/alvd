package runner

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/daemon"
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
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	d, err := daemon.New(r.cfg)
	if err != nil {
		return err
	}

	ech := d.Start(ctx)
	defer d.Close()

	for {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
			return nil
		case err := <-ech:
			log.Errorf("error: %s", err)
		}
	}
}
