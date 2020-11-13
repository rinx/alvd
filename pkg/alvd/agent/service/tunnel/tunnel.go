package tunnel

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/rancher/remotedialer"
	"github.com/rinx/alvd/internal/log"
)

type tunnel struct {
	*Config
	cancel context.CancelFunc
}

type Config struct {
	ServerAddress string

	AgentName string
	AgentPort uint
}

type Tunnel interface {
	Close()
}

func Connect(ctx context.Context, cfg *Config) (Tunnel, <-chan error) {
	ctx, cancel := context.WithCancel(ctx)
	ech := make(chan error, 1)

	headers := http.Header{
		"X-ALVD-ID":        []string{cfg.AgentName},
		"X-ALVD-GRPC-PORT": []string{strconv.Itoa(int(cfg.AgentPort))},
	}

	go func() {
		defer close(ech)
		for {
			remotedialer.ClientConnect(
				ctx,
				fmt.Sprintf("ws://%s/connect", cfg.ServerAddress),
				headers,
				nil,
				connectAuthorizer,
				onConnectFunc(cfg.ServerAddress),
			)

			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err != nil {
					log.Errorf("error: %s", err)
				}
				return
			default:
			}
		}
	}()

	return &tunnel{
		Config: cfg,
		cancel: cancel,
	}, ech
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
