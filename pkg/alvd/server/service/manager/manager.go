package manager

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/rinx/alvd/internal/errors"
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

	clients []Client
	mu      sync.RWMutex
}

type Client struct {
	Key  string
	Port int
	Addr string

	StoredIndex      int
	UncommittedIndex int
	IsIndexing       bool
}

type Manager interface {
	Start(ctx context.Context) <-chan error
	Close() error
	GetClient(addr string) (vald.Client, error)
	GetAgentClient(addr string) (core.AgentClient, error)
	GetClientsList() []Client

	GetAgentCount() int

	Broadcast(ctx context.Context, f func(ctx context.Context, client vald.Client) error) error
	Range(ctx context.Context, concurrency int, f func(ctx context.Context, client vald.Client) error) error
}

func New(
	tun tunnel.Tunnel,
	checkIndexInterval time.Duration,
) (Manager, error) {
	return &manager{
		interval: checkIndexInterval,
		tunnel:   tun,
		clients:  make([]Client, 0),
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

func (m *manager) GetAgentCount() int {
	return len(m.GetClientsList())
}

func (m *manager) GetClientsList() []Client {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.clients
}

func (m *manager) updateClientsList(ctx context.Context, ech chan error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmap := m.tunnel.Clients()

	m.clients = make([]Client, 0, len(cmap))

	log.Debugf("%d clients are found. updating client list.", len(cmap))

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

		m.clients = append(m.clients, Client{
			Key:              key,
			Port:             port,
			Addr:             addr,
			StoredIndex:      int(idxInfo.GetStored()),
			UncommittedIndex: int(idxInfo.GetUncommitted()),
			IsIndexing:       idxInfo.GetIndexing(),
		})
	}

	sort.Slice(m.clients, func(i, j int) bool {
		ix := m.clients[i].StoredIndex + m.clients[i].UncommittedIndex
		jx := m.clients[j].StoredIndex + m.clients[j].UncommittedIndex
		return ix < jx
	})

	log.Debugf("client list updated: %#v", m.clients)
}

func collectError(ctx context.Context, err *error, ech <-chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-ech:
			if e != nil {
				*err = errors.Wrap(*err, e.Error())
			}
		}
	}
}

func (m *manager) Broadcast(
	ctx context.Context,
	f func(ctx context.Context, client vald.Client) error,
) (err error) {
	cl := m.GetClientsList()
	wg := sync.WaitGroup{}

	ech := make(chan error, 1)
	defer close(ech)

	ectx, cancel := context.WithCancel(ctx)
	go collectError(ectx, &err, ech)

	for _, c := range cl {
		wg.Add(1)

		go func(c Client) {
			defer wg.Done()

			client, err := m.GetClient(c.Addr)
			if err != nil {
				ech <- errors.Errorf("Cannot get client for %s: %s", c.Addr, err)
				return
			}

			err = f(ctx, client)
			if err != nil {
				ech <- errors.Errorf("Error from %s: %s", c.Addr, err)
			}
		}(c)
	}

	wg.Wait()

	cancel()

	return err
}

func (m *manager) Range(
	ctx context.Context,
	concurrency int,
	f func(ctx context.Context, client vald.Client) error,
) (err error) {
	cl := m.GetClientsList()
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrency)

	ech := make(chan error, 1)
	defer close(ech)

	ectx, cancel := context.WithCancel(ctx)
	go collectError(ectx, &err, ech)

	for _, c := range cl {
		wg.Add(1)

		go func(c Client) {
			defer wg.Done()

			defer func() {
				select {
				case <-ctx.Done():
				case <-semaphore:
				}
			}()

			select {
			case <-ctx.Done():
				return
			case semaphore <- struct{}{}:
			}

			client, err := m.GetClient(c.Addr)
			if err != nil {
				ech <- errors.Errorf("Cannot get client for %s: %s", c.Addr, err)
				return
			}

			err = f(ctx, client)
			if err != nil {
				ech <- errors.Errorf("Error from %s: %s", c.Addr, err)
			}
		}(c)
	}

	wg.Wait()

	cancel()

	return err
}
