package tunnel

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rancher/remotedialer"
	"github.com/rinx/alvd/internal/net/tcp"
)

type tunnel struct {
	server *remotedialer.Server
}

type Tunnel interface {
	Handler() http.Handler
	AddPeer(url, id, token string)
	Dialer(key string, timeout time.Duration) tcp.Dialer
}

func New() (Tunnel, error) {
	return &tunnel{
		server: remotedialer.New(
			authorizer,
			remotedialer.DefaultErrorWriter,
		),
	}, nil
}

func authorizer(req *http.Request) (clientKey string, authed bool, err error) {
	return "", true, nil
}

func (t *tunnel) Handler() http.Handler {
	return t.server
}

func (t *tunnel) AddPeer(url, id, token string) {
	t.server.AddPeer(url, id, token)
}

func (t *tunnel) Dialer(key string, timeout time.Duration) tcp.Dialer {
	return &dialer{
		key:     key,
		timeout: timeout,
		dialer:  t.server.Dialer(key, timeout),
	}
}

type dialer struct {
	key     string
	timeout time.Duration
	dialer  remotedialer.Dialer
}

func (d *dialer) GetDialer() func(ctx context.Context, network, address string) (net.Conn, error) {
	return d.DialContext
}

func (d *dialer) StartDialerCache(ctx context.Context) {
}

func (d *dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := d.dialer(network, address)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
