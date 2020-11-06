package runner

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rinx/alvd/internal/log"
	"github.com/rinx/alvd/pkg/alvd/server/config"
	"github.com/rinx/alvd/pkg/alvd/server/tunnel"
)

type runner struct {
	cfg *config.Config
}

type Runner interface {
	Start(ctx context.Context) error
}

func New(cfg *config.Config) (Runner, error) {
	return &runner{
		cfg: cfg,
	}, nil
}

func (r *runner) Start(ctx context.Context) error {
	tun, err := tunnel.New()
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.Handle("/connect", tun.Handler())

	log.Infof("listen: %s", r.cfg.Addr)
	return http.ListenAndServe(r.cfg.Addr, router)
}
