package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/daemon"
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

	wg := sync.WaitGroup{}

	for {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
			wg.Wait()
			return nil
		case err := <-ech:
			if err != context.Canceled {
				log.Errorf("error: %s", err)
			}
		}
	}
}
