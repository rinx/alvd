package manager

import (
	"context"
	"fmt"
	"sort"
	"sync"
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

	clients []client
	mu      sync.RWMutex
}

type client struct {
	key  string
	port int
	addr string

	indexInfo *payload.Info_Index_Count
}

type Manager interface {
	Start(ctx context.Context) <-chan error
	Close() error
	GetClient(addr string) (vald.Client, error)
	GetAgentClient(addr string) (core.AgentClient, error)

	Broadcast(ctx context.Context, f func(ctx context.Context, client vald.Client) error) error
}

func New(tun tunnel.Tunnel) (Manager, error) {
	return &manager{
		interval: 5000 * time.Millisecond,
		tunnel:   tun,
		clients:  make([]client, 0),
	}, nil
}

func (m *manager) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 1)

	go func() {
		defer close(ech)

		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			m.updateClientsList(ctx, ech)

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

func (m *manager) GetClient(addr string) (vald.Client, error) {
	conn, err := m.getConn(addr)
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

func (m *manager) getClientsList() []client {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.clients
}

func (m *manager) updateClientsList(ctx context.Context, ech chan error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmap := m.tunnel.Clients()

	m.clients = make([]client, 0, len(cmap))

	for key, port := range cmap {
		addr := toAddr(key, port)

		cli, err := m.GetAgentClient(addr)
		if err != nil {
			ech <- err
			continue
		}

		idxInfo, err := cli.IndexInfo(ctx, &payload.Empty{})
		if err != nil {
			ech <- err
			continue
		}

		m.clients = append(m.clients, client{
			key:       key,
			port:      port,
			addr:      addr,
			indexInfo: idxInfo,
		})
	}

	sort.Slice(m.clients, func(i, j int) bool {
		ix := m.clients[i].indexInfo.Stored + m.clients[i].indexInfo.Uncommitted
		jx := m.clients[j].indexInfo.Stored + m.clients[j].indexInfo.Uncommitted
		return ix < jx
	})
}

func (m *manager) Broadcast(
	ctx context.Context,
	f func(ctx context.Context, client vald.Client) error,
) error {
	cl := m.getClientsList()
	wg := sync.WaitGroup{}

	for _, c := range cl {
		wg.Add(1)

		go func(c client) {
			defer wg.Done()

			client, err := m.GetClient(c.addr)
			if err != nil {
				log.Errorf("%s", err)
				return
			}

			err = f(ctx, client)
			if err != nil {
				log.Errorf("%s", err)
			}
		}(c)
	}

	wg.Wait()

	return nil
}
