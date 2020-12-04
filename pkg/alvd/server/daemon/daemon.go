package daemon

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/service/gateway"
	"github.com/rinx/alvd/pkg/alvd/server/service/gateway/handler"
	"github.com/rinx/alvd/pkg/alvd/server/service/indexer"
	"github.com/rinx/alvd/pkg/alvd/server/service/manager"
	"github.com/rinx/alvd/pkg/alvd/server/service/tunnel"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

type daemon struct {
	addr     string
	grpcAddr string

	cancel context.CancelFunc

	gateway gateway.Gateway
	tunnel  tunnel.Tunnel
	manager manager.Manager
	handler vald.Server
	indexer indexer.Indexer
}

type Daemon interface {
	Start(ctx context.Context) <-chan error
	Close() error
}

func New(cfg *config.Config) (Daemon, error) {
	tun, err := tunnel.New()
	if err != nil {
		return nil, err
	}

	m, err := manager.New(tun, cfg.CheckIndexInterval)
	if err != nil {
		return nil, err
	}

	i, err := indexer.New(
		m,
		cfg.CheckIndexInterval,
		cfg.CreateIndexThreshold,
	)
	if err != nil {
		return nil, err
	}

	h := handler.New(m, cfg.Replicas)

	g, err := gateway.New(h, cfg.GRPCHost, cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &daemon{
		addr:    cfg.Addr,
		gateway: g,
		tunnel:  tun,
		manager: m,
		handler: h,
		indexer: i,
	}, nil
}

func (d *daemon) Start(ctx context.Context) <-chan error {
	ctx, d.cancel = context.WithCancel(ctx)

	sech := d.startHTTPServer(ctx)
	mech := d.manager.Start(ctx)
	iech := d.indexer.Start(ctx)
	gech := d.gateway.Start(ctx)
	ech := make(chan error, 1)

	go func() {
		defer close(ech)

		for {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			case err := <-sech:
				ech <- err
			case err := <-mech:
				ech <- err
			case err := <-iech:
				ech <- err
			case err := <-gech:
				ech <- err
			}
		}
	}()

	return ech
}

func (d *daemon) startHTTPServer(ctx context.Context) <-chan error {
	ech := make(chan error, 1)
	router := mux.NewRouter()
	router.Handle("/connect", d.tunnel.Handler())

	go func() {
		defer close(ech)

		for {
			log.Infof("websocket server starting on %s", d.addr)
			err := http.ListenAndServe(d.addr, router)
			if err != nil {
				ech <- err
			}

			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			default:
			}
		}
	}()

	return ech
}

func (d *daemon) Close() (err error) {
	d.gateway.Close()
	d.indexer.Close()
	d.manager.Close()

	d.cancel()

	return nil
}
