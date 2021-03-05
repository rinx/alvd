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
	agentName string
	agentPort int

	cancel       context.CancelFunc
	cancelByAddr map[string]context.CancelFunc

	connectCh    chan string
	disconnectCh chan string
}

type Tunnel interface {
	Start(ctx context.Context) <-chan error
	Connect(addr string)
	Disconnect(addr string)
	Close()
}

func New(name string, port int) Tunnel {
	return &tunnel{
		agentName:    name,
		agentPort:    port,
		cancelByAddr: make(map[string]context.CancelFunc, 0),
	}
}

func (t *tunnel) Start(ctx context.Context) <-chan error {
	ctx, t.cancel = context.WithCancel(ctx)
	ech := make(chan error, 1)
	t.connectCh = make(chan string, 10)
	t.disconnectCh = make(chan string, 10)

	go func() {
		defer close(ech)
		defer close(t.connectCh)
		defer close(t.disconnectCh)
		var err error
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				if err != nil && err != context.Canceled {
					log.Errorf("error: %s", err)
				}
				return
			case addr := <-t.connectCh:
				host, port, err := net.SplitHostPort(addr)
				if err == nil {
					ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
					if err == nil {
						for _, ip := range ips {
							t.connect(ctx, net.JoinHostPort(ip.String(), port))
						}
						continue
					}
				}

				t.connect(ctx, addr)
			case addr := <-t.disconnectCh:
				t.disconnect(ctx, addr)
			}
		}
	}()

	return ech
}

func (t *tunnel) connect(ctx context.Context, addr string) {
	ctx, t.cancelByAddr[addr] = context.WithCancel(ctx)

	headers := http.Header{
		"X-ALVD-ID":        []string{t.agentName},
		"X-ALVD-GRPC-PORT": []string{strconv.Itoa(t.agentPort)},
	}

	go func() {
		for {
			remotedialer.ClientConnect(
				ctx,
				fmt.Sprintf("ws://%s/connect", addr),
				headers,
				nil,
				connectAuthorizer,
				onConnectFunc(addr),
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
}

func (t *tunnel) disconnect(ctx context.Context, addr string) {
	cancel, ok := t.cancelByAddr[addr]
	if ok {
		cancel()
		delete(t.cancelByAddr, addr)
	}
}

func (t *tunnel) Connect(addr string) {
	t.connectCh <- addr
}

func (t *tunnel) Disconnect(addr string) {
	t.disconnectCh <- addr
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
