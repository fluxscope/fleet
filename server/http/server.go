package http

import (
	"context"
	"github.com/fluxscope/fleet"
	"github.com/fluxscope/fleet/pkg/log"
	"log/slog"
	"net"
	"net/http"
	"sync"
)

var _ fleet.Service = new(Server)

type Server struct {
	*http.Server
	addr         string
	listener     net.Listener
	logger       *slog.Logger
	shutdownOnce sync.Once
}

func (s *Server) ID() string {
	return "http"
}

type Option func(*Server)

func WithHTTPServer(hs *http.Server) Option {
	return func(s *Server) {
		s.Server = hs
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithListener(l net.Listener) Option {
	return func(s *Server) {
		s.listener = l
	}
}

func NewServer(opts ...Option) *Server {
	s := new(Server)
	for _, opt := range opts {
		opt(s)
	}

	if s.Server == nil {
		s.Server = &http.Server{Addr: s.addr}
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	if s.logger == nil {
		s.logger = log.FromContext(ctx).With("service", s.ID())
	}

	if s.listener == nil {
		s.logger.Info("starting HTTP server", "addr", s.Server.Addr)
		return s.Server.ListenAndServe()
	} else {
		return s.Serve(s.listener)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	var err error
	s.shutdownOnce.Do(func() {
		if s.Server != nil {
			s.logger.Info("shutdown HTTP server")
			err = s.Server.Shutdown(ctx)
		}
	})
	return err
}
