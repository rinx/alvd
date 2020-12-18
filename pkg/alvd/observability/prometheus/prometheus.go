package prometheus

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
)

type prom struct {
	exporter *prometheus.Exporter
}

type Prometheus interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func New() (Prometheus, error) {
	conf := prometheus.Config{}

	exporter, err := prometheus.NewExportPipeline(conf)
	if err != nil {
		return nil, err
	}

	otel.SetMeterProvider(exporter.MeterProvider())

	return &prom{
		exporter: exporter,
	}, nil
}

func (p *prom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.exporter.ServeHTTP(w, r)
}
