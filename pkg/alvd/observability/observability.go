package observability

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/observability/metrics"
	"github.com/rinx/alvd/pkg/alvd/observability/prometheus"
)

type obs struct {
	prometheus       prometheus.Prometheus
	metricServerAddr string
}

type Obs interface {
	Start(ctx context.Context) error
}

func New(cfg *Config) (Obs, error) {
	addr := fmt.Sprintf("%s:%d", cfg.MetricsHost, cfg.MetricsPort)

	prom, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	return &obs{
		prometheus:       prom,
		metricServerAddr: addr,
	}, nil
}

func (o *obs) Start(ctx context.Context) (err error) {
	metrics.Init()

	sech := o.StartMetricServer(ctx)
	mech := metrics.GetMeter().Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			case err = <-sech:
				log.Errorf("error: %s", err)
			case err = <-mech:
				log.Errorf("error: %s", err)
			}
		}
	}()

	return nil
}

func (o *obs) StartMetricServer(ctx context.Context) <-chan error {
	ech := make(chan error, 1)
	router := mux.NewRouter()
	router.HandleFunc("/metrics", o.prometheus.ServeHTTP)

	go func() {
		defer close(ech)

		for {
			log.Infof("metrics server starting on %s", o.metricServerAddr)
			err := http.ListenAndServe(o.metricServerAddr, router)
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
