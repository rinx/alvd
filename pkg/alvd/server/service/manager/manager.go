package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/server/service/tunnel"
	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"google.golang.org/grpc"
)

type manager struct {
	interval time.Duration

	tunnel tunnel.Tunnel
}

type Manager interface {
	Start(ctx context.Context) <-chan error
	Close() error
	GetClient() (vald.Client, error)
}

func New(tun tunnel.Tunnel) (Manager, error) {
	return &manager{
		interval: 5000 * time.Millisecond,
		tunnel:   tun,
	}, nil
}

func (m *manager) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	go func() {
		defer close(ech)

		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			for key, port := range m.tunnel.Clients() {
				log.Infof("id: %s, port: %d", key, port)

				client, err := m.GetAgentClient(toAddr(key, port))
				if err != nil {
					ech <- err
					continue
				}

				res, err := client.IndexInfo(ctx, &payload.Empty{})
				if err != nil {
					ech <- err
					continue
				}

				log.Infof("%#v", res)
			}

			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil {
					log.Errorf("error: %s", err)
				}
				return
			case <-ticker.C:
			}
		}
	}()

	return ech
}

func (m *manager) Close() error {
	return nil
}

func toAddr(key string, port int) string {
	return fmt.Sprintf("%s:%d", key, port)
}

func (m *manager) getConn(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithContextDialer(m.tunnel.ContextDialer),
	)
}

func (m *manager) GetClient() (vald.Client, error) {
	conn, err := m.getConn("agent:8081")
	if err != nil {
		return nil, err
	}

	return vald.NewValdClient(conn), nil
}

func (m *manager) GetAgentClient(addr string) (core.AgentClient, error) {
	conn, err := m.getConn(addr)
	if err != nil {
		return nil, err
	}

	return core.NewAgentClient(conn), nil
}
