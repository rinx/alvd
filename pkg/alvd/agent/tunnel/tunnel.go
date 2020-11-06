package tunnel

import (
	"context"
	"fmt"
	"net"

	"github.com/rancher/remotedialer"
)

type tunnel struct {
	address string
	cancel  context.CancelFunc
}

type Tunnel interface {
	Close()
}

func Connect(ctx context.Context, address string) Tunnel {
	ctx, cancel := context.WithCancel(ctx)

	if address != "" {
		remotedialer.ClientConnect(
			ctx,
			fmt.Sprintf("ws://%s/connect", address),
			nil,
			nil,
			connectAuthorizer,
			onConnect,
		)
	}

	return &tunnel{
		address: address,
		cancel:  cancel,
	}
}

func (t *tunnel) Close() {
	t.cancel()
}

func connectAuthorizer(proto, address string) bool {
	host, _, err := net.SplitHostPort(address)
	return err == nil && proto == "tcp" && host == "127.0.0.1"
}

func onConnect(ctx context.Context) error {
	return nil
}
