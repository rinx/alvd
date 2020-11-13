package gateway

import (
	"context"
	"fmt"
	"net"

	"github.com/rinx/alvd/internal/log"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"google.golang.org/grpc"
)

type gateway struct {
	handler vald.Server
}

type Gateway interface {
	Start(ctx context.Context) <-chan error
	Close() error
}

func New(handler vald.Server) (Gateway, error) {

	return &gateway{
		handler: handler,
	}, nil
}

func (g *gateway) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	gech := g.startGRPCServer(ctx)

	go func() {
		defer close(ech)

		var err error
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil {
					log.Errorf("error: %s", err)
				}
				return
			case err = <-gech:
				ech <- err
			}
		}
	}()

	return ech
}

func (g *gateway) startGRPCServer(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	server := grpc.NewServer()
	vald.RegisterValdServer(server, g.handler)

	go func() {
		defer close(ech)

		for {
			addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8082)

			log.Infof("listen: %s", addr)

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				ech <- err
			} else {
				err = server.Serve(lis)
				if err != nil {
					ech <- err
				}
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

	return ech
}

func (g *gateway) Close() error {
	return nil
}
