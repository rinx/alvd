package daemon

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/tunnel"
)

type daemon struct {
	addr string

	cancel context.CancelFunc

	tunnel tunnel.Tunnel
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

	return &daemon{
		addr:   cfg.Addr,
		tunnel: tun,
	}, nil
}

func (d *daemon) Start(ctx context.Context) <-chan error {
	ctx, d.cancel = context.WithCancel(ctx)

	sech := make(chan error, 1)
	router := mux.NewRouter()
	router.Handle("/connect", d.tunnel.Handler())

	go func() {
		defer close(sech)

		for {
			log.Infof("listen: %s", d.addr)
			err := http.ListenAndServe(d.addr, router)
			if err != nil {
				sech <- err
			}

			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil {
					log.Errorf("error: %s", err)
				}
				return
			default:
			}
		}
	}()

	ech := make(chan error, 1)

	go func() {
		defer close(ech)

		for {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil {
					log.Errorf("error: %s", err)
				}
				return
			case err := <-sech:
				ech <- err
			}
		}
	}()

	return ech
}

func (d *daemon) Close() error {
	d.cancel()

	return nil
}
