package daemon

import (
	"context"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/service/tunnel"
)

type daemon struct {
	serverAddress string

	agentName string
	agentPort uint

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
		agentName:     cfg.AgentName,
		agentPort:     cfg.AgentPort,
	}, nil
}

func (d *daemon) Start(ctx context.Context) error {
	ctx, d.cancel = context.WithCancel(ctx)

	var tunEch <-chan error
	d.tunnel, tunEch = tunnel.Connect(ctx, &tunnel.Config{
		ServerAddress: d.serverAddress,
		AgentName:     d.agentName,
		AgentPort:     d.agentPort,
	})

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
