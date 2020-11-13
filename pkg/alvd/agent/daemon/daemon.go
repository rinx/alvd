package daemon

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/tunnel"
)

type daemon struct {
	serverAddress string

	cancel context.CancelFunc

	tunnel tunnel.Tunnel
}

type Daemon interface {
	Start(ctx context.Context) error
	Close() error
}

func New(cfg *config.Config) (Daemon, error) {

	return &daemon{
		serverAddress: cfg.ServerAddress,
	}, nil
}

func (d *daemon) Start(ctx context.Context) error {
	ctx, d.cancel = context.WithCancel(ctx)

	var tunEch <-chan error
	d.tunnel, tunEch = tunnel.Connect(ctx, d.serverAddress)

	go func() {
		var err error
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil {
					log.Errorf("error: %s", err)
				}
				return
			case err = <-tunEch:
				log.Errorf("error: %s", err)
			}
		}
	}()

	return nil
}

func (d *daemon) Close() error {
	defer d.cancel()
	defer d.tunnel.Close()

	return nil
}
