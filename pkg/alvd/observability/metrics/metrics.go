package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/rinx/alvd/internal/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

var (
	once     sync.Once
	instance *meter
)

type meter struct {
	meter metric.Meter

	mch chan metric.Measurement

	interval time.Duration

	ls []attribute.KeyValue
	ms []metric.Measurement

	collectorsMu sync.RWMutex
	collectors   []func() (metric.Measurement, error)
}

type Meter interface {
	Start(ctx context.Context) <-chan error
	RegisterCollectors(collectors ...func() (metric.Measurement, error))
	Sink(measure metric.Measurement)
	Meter() metric.Meter
}

func Init(interval time.Duration) {
	once.Do(func() {
		instance = &meter{
			meter:    global.Meter("rinx.github.io/alvd"),
			interval: interval,
			mch:      make(chan metric.Measurement, 10),
		}
	})
}

func GetMeter() Meter {
	return instance
}

func (m *meter) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	var measure metric.Measurement

	go func() {
		defer close(ech)

		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			case <-ticker.C:
				m.collect(ech)

				m.meter.RecordBatch(ctx, m.ls, m.ms...)
				m.ms = make([]metric.Measurement, 0)
			case measure = <-m.mch:
				m.ms = append(m.ms, measure)
			}
		}
	}()

	return ech
}

func (m *meter) collect(ech chan error) {
	m.collectorsMu.RLock()
	defer m.collectorsMu.RUnlock()

	for _, collector := range m.collectors {
		measure, err := collector()
		if err != nil {
			ech <- err
			continue
		}

		m.ms = append(m.ms, measure)
	}
}

func (m *meter) RegisterCollectors(collectors ...func() (metric.Measurement, error)) {
	m.collectorsMu.Lock()
	defer m.collectorsMu.Unlock()

	m.collectors = append(m.collectors, collectors...)
}

func (m *meter) Sink(measure metric.Measurement) {
	m.mch <- measure
}

func (m *meter) Meter() metric.Meter {
	return m.meter
}
