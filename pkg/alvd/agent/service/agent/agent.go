package agent

import (
	"context"
	"fmt"
	"net"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/agent/service/agent/handler"
	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"google.golang.org/grpc"
)

type agent struct {
	handler handler.Server
	addr    string
}

type Agent interface {
	Start(ctx context.Context) <-chan error
	Close() error
}

func New(handler handler.Server, host string, port int) (Agent, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	return &agent{
		handler: handler,
		addr:    addr,
	}, nil
}

func (a *agent) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	gech := a.startGRPCServer(ctx)

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
			case err = <-gech:
				ech <- err
			}
		}
	}()

	return ech

}

func (a *agent) startGRPCServer(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	server := grpc.NewServer()
	core.RegisterAgentServer(server, a.handler)
	vald.RegisterValdServer(server, a.handler)

	go func() {
		defer close(ech)

		for {
			log.Infof("agent gRPC API starting on %s", a.addr)

			lis, err := net.Listen("tcp", a.addr)
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

func (a *agent) Close() error {
	return nil
}
