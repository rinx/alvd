package indexer

import (
	"context"
	"time"

	"github.com/rinx/alvd/internal/errors"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/internal/net/grpc/codes"
	"github.com/rinx/alvd/internal/net/grpc/status"
	"github.com/rinx/alvd/pkg/alvd/server/service/manager"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

type indexer struct {
	manager manager.Manager

	interval  time.Duration
	threshold int
}

type Indexer interface {
	Start(ctx context.Context) <-chan error
	Close() error
}

func New(
	manager manager.Manager,
	checkIndexInterval time.Duration,
	createIndexThreshold int,
) (Indexer, error) {
	return &indexer{
		manager:   manager,
		interval:  checkIndexInterval,
		threshold: createIndexThreshold,
	}, nil
}

func (i *indexer) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	go func() {
		defer close(ech)

		ticker := time.NewTicker(i.interval)
		defer ticker.Stop()

		for {
			i.checkIndex(ctx, ech)

			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			case <-ticker.C:
			}
		}
	}()

	return ech
}

func (i *indexer) checkIndex(ctx context.Context, ech chan error) {
	cl := i.manager.GetClientsList()

	for _, c := range cl {
		if c.UncommittedIndex >= i.threshold {
			client, err := i.manager.GetAgentClient(c.Addr)
			if err != nil {
				ech <- errors.Errorf("Cannot get client for %s: %s", c.Addr, err)
				continue
			}

			log.Debugf("create index started for %s", c.Addr)

			_, err = client.CreateIndex(ctx, &payload.Control_CreateIndexRequest{})
			if err != nil {
				if status.Code(err) != codes.FailedPrecondition {
					ech <- errors.Errorf("Error occurred when creating index for %s: %s", c.Addr, err)
				}
			}

			log.Debugf("create index finished for %s", c.Addr)
		}
	}
}

func (i *indexer) Close() error {
	return nil
}
