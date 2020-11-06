package tunnel

import (
	"context"
	"fmt"
	"net"

	"github.com/rancher/remotedialer"
	"github.com/rinx/alvd/internal/log"
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
		go func() {
			for {
				remotedialer.ClientConnect(
					ctx,
					fmt.Sprintf("ws://%s/connect", address),
					nil,
					nil,
					connectAuthorizer,
					onConnectFunc(address),
				)

				select {
				case <-ctx.Done():
					err := ctx.Err()
					if err != nil {
						log.Errorf("%s", err)
						return
					}
				default:
				}
			}
		}()
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

func onConnectFunc(address string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		log.Infof("connected to: %s", address)
		return nil
	}
}
