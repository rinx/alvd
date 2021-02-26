package daemon

import (
	"context"

	"github.com/rinx/alvd/internal/errgroup"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/agent/config"
	"github.com/rinx/alvd/pkg/alvd/agent/service/agent"
	"github.com/rinx/alvd/pkg/alvd/agent/service/agent/handler"
	"github.com/rinx/alvd/pkg/alvd/agent/service/tunnel"
	"github.com/rinx/alvd/pkg/alvd/observability/metrics"
	"github.com/rinx/alvd/pkg/vald/agent/ngt/service"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/unit"
)

type daemon struct {
	serverAddresses []string

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
		serverAddresses: cfg.ServerAddresses,
		agentName:       cfg.AgentName,
		grpcPort:        cfg.GRPCPort,
		agent:           a,
		ngt:             ngt,
		handler:         h,
	}, nil
}

func (d *daemon) Start(ctx context.Context) <-chan error {
	ctx, d.cancel = context.WithCancel(ctx)

	d.tunnel = tunnel.New(d.agentName, d.grpcPort)
	tunEch := d.tunnel.Start(ctx)
	for _, addr := range d.serverAddresses {
		d.tunnel.Connect(addr)
	}

	nech := d.ngt.Start(ctx)
	gech := d.agent.Start(ctx)

	ech := make(chan error, 1)

	go func() {
		defer close(ech)
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

	err := d.registerMetrics()
	if err != nil {
		ech <- err
	}

	return ech
}

func (d *daemon) Close() error {
	defer d.tunnel.Close()

	defer d.cancel()

	return nil
}

func (d *daemon) registerMetrics() (err error) {
	meter := metrics.GetMeter()

	_, err = meter.Meter().NewInt64UpDownSumObserver(
		"rinx.github.io/alvd/agent/ngt/stored",
		func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(int64(d.ngt.Len()))
		},
		metric.WithDescription("NGT number of stored indices"),
		metric.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = meter.Meter().NewInt64UpDownSumObserver(
		"rinx.github.io/alvd/agent/ngt/uncommitted",
		func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(int64(d.ngt.InsertVCacheLen() + d.ngt.DeleteVCacheLen()))
		},
		metric.WithDescription("NGT number of uncommitted indices"),
		metric.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return err
	}

	return nil
}
