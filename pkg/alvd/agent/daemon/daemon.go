package daemon

import (
	"context"

	"github.com/rinx/alvd/internal/errgroup"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/service/agent"
	"github.com/rinx/alvd/pkg/alvd/agent/service/agent/handler"
	"github.com/rinx/alvd/pkg/alvd/agent/service/tunnel"
	"github.com/rinx/alvd/pkg/vald/agent/ngt/service"
)

type daemon struct {
	serverAddress string

	agentName string
	grpcPort  int

	cancel context.CancelFunc

	agent   agent.Agent
	tunnel  tunnel.Tunnel
	ngt     service.NGT
	handler handler.Server
}

type Daemon interface {
	Start(ctx context.Context) <-chan error
	Close() error
}

func New(cfg *config.Config) (Daemon, error) {
	ngt, err := service.New(
		cfg.NGTConfig,
		service.WithErrGroup(errgroup.Get()),
		service.WithEnableInMemoryMode(cfg.NGTConfig.EnableInMemoryMode),
		service.WithIndexPath(cfg.NGTConfig.IndexPath),
		service.WithAutoIndexCheckDuration(cfg.NGTConfig.AutoIndexCheckDuration),
		service.WithAutoIndexDurationLimit(cfg.NGTConfig.AutoIndexDurationLimit),
		service.WithAutoSaveIndexDuration(cfg.NGTConfig.AutoSaveIndexDuration),
		service.WithAutoIndexLength(cfg.NGTConfig.AutoIndexLength),
		service.WithInitialDelayMaxDuration(cfg.NGTConfig.InitialDelayMaxDuration),
		service.WithMinLoadIndexTimeout(cfg.NGTConfig.MinLoadIndexTimeout),
		service.WithMaxLoadIndexTimeout(cfg.NGTConfig.MaxLoadIndexTimeout),
		service.WithLoadIndexTimeoutFactor(cfg.NGTConfig.LoadIndexTimeoutFactor),
		service.WithDefaultPoolSize(cfg.NGTConfig.DefaultPoolSize),
		service.WithDefaultRadius(cfg.NGTConfig.DefaultRadius),
		service.WithDefaultEpsilon(cfg.NGTConfig.DefaultEpsilon),
		service.WithProactiveGC(cfg.NGTConfig.EnableProactiveGC),
	)
	if err != nil {
		return nil, err
	}

	h := handler.New(cfg.AgentName, ngt)

	a, err := agent.New(h, cfg.GRPCHost, cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &daemon{
		serverAddress: cfg.ServerAddress,
		agentName:     cfg.AgentName,
		grpcPort:      cfg.GRPCPort,
		agent:         a,
		ngt:           ngt,
		handler:       h,
	}, nil
}

func (d *daemon) Start(ctx context.Context) <-chan error {
	ctx, d.cancel = context.WithCancel(ctx)

	var tunEch <-chan error
	d.tunnel, tunEch = tunnel.Connect(ctx, &tunnel.Config{
		ServerAddress: d.serverAddress,
		AgentName:     d.agentName,
		AgentPort:     d.grpcPort,
	})

	nech := d.ngt.Start(ctx)
	gech := d.agent.Start(ctx)

	ech := make(chan error, 1)

	go func() {
		var err error
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			case err = <-tunEch:
				ech <- err
			case err = <-nech:
				ech <- err
			case err = <-gech:
				ech <- err
			}
		}
	}()

	return ech
}

func (d *daemon) Close() error {
	defer d.tunnel.Close()

	defer d.cancel()

	return nil
}
