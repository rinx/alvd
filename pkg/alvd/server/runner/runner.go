package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/daemon"
)

type runner struct {
	cfg *config.Config
}

type Runner interface {
	Start(ctx context.Context) error
}

func New(cfg *config.Config) (_ Runner, err error) {
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

	wg := sync.WaitGroup{}

	if r.cfg.AgentEnabled {
		wg.Add(1)

		go func() {
			defer wg.Done()
			agent.Run(r.cfg.AgentOpts)
		}()
	}

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
