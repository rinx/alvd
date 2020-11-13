package tunnel

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rancher/remotedialer"
	"github.com/rinx/alvd/internal/errors"
)

type tunnel struct {
	server *remotedialer.Server

	// map of id-port
	clients map[string]int
	mu      sync.Mutex

	dialerDefaultHost    string
	dialerDefaultProto   string
	dialerDefaultTimeout time.Duration
}

type Tunnel interface {
	Handler() http.Handler
	RemovePeer(key string)
	HasSession(key string) error
	Clients() map[string]int
	ContextDialer(ctx context.Context, address string) (net.Conn, error)
}

func New() (Tunnel, error) {
	t := &tunnel{
		clients:              make(map[string]int, 0),
		dialerDefaultHost:    "127.0.0.1",
		dialerDefaultProto:   "tcp",
		dialerDefaultTimeout: 15 * time.Second,
	}

	t.server = remotedialer.New(
		t.authorizer(),
		remotedialer.DefaultErrorWriter,
	)

	return t, nil
}

func (t *tunnel) authorizer() func(req *http.Request) (clientKey string, authed bool, err error) {
	return func(req *http.Request) (clientKey string, authed bool, err error) {
		name := req.Header.Get("X-ALVD-ID")
		port := req.Header.Get("X-ALVD-GRPC-PORT")

		iport, err := strconv.Atoi(port)
		if err != nil {
			return "", false, err
		}

		t.mu.Lock()
		defer t.mu.Unlock()

		_, ok := t.clients[name]
		if ok {
			return "", false, errors.Errorf("same client name %s already exists", name)
		}

		t.clients[name] = iport

		return name, true, nil
	}
}

func (t *tunnel) Handler() http.Handler {
	return t.server
}

func (t *tunnel) RemovePeer(key string) {
	t.server.RemovePeer(key)

	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.clients, key)
}

func (t *tunnel) HasSession(key string) error {
	if t.server.HasSession(key) {
		return nil
	}

	t.RemovePeer(key)

	return errors.Errorf("session has closed for client: %s", key)
}

func (t *tunnel) Clients() map[string]int {
	return t.clients
}

func (t *tunnel) ContextDialer(ctx context.Context, address string) (net.Conn, error) {
	key, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	if err = t.HasSession(key); err != nil {
		return nil, err
	}

	return t.server.Dial(
		key,
		t.dialerDefaultTimeout,
		t.dialerDefaultProto,
		fmt.Sprintf("%s:%s", t.dialerDefaultHost, port),
	)
}
