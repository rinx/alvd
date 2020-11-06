package tunnel

import (
	"net/http"
	"time"

	"github.com/rancher/remotedialer"
)

type dialer struct {
	server *remotedialer.Server
}

type Dialer interface {
	Handler() http.Handler
	Dialer(key string, timeout time.Duration) remotedialer.Dialer
}

func New() (Dialer, error) {
	return &dialer{
		server: remotedialer.New(
			authorizer,
			remotedialer.DefaultErrorWriter,
		),
	}, nil
}

func (d *dialer) Handler() http.Handler {
	return d.server
}

func (d *dialer) Dialer(key string, timeout time.Duration) remotedialer.Dialer {
	return d.server.Dialer(key, timeout)
}

func authorizer(req *http.Request) (clientKey string, authed bool, err error) {
	return "", true, nil
}
