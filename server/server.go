package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegodesousas/httpserver/middlewares"
)

var DefaultMonitorWrapper = func(path string, handler http.Handler) http.Handler {
	return handler
}

type (
	MonitorWrapper func(path string, handler http.Handler) http.Handler
)

type Server struct {
	http.Server
	routes         []Route
	router         Router
	monitorWrapper MonitorWrapper
}

type ShutdownHandler func(ctx context.Context) error

func New(configs ...Option) *Server {
	server := &Server{
		Server: http.Server{
			Addr: ":8080",
		},
		monitorWrapper: DefaultMonitorWrapper,
	}

	server.Server.Handler = middlewares.Middlewares(
		server,
		middlewares.PanicRecoveryMiddleware,
	)

	for _, config := range configs {
		config(server)
	}

	for _, r := range server.routes {
		server.router.Resource(r.Method, r.Path, server.monitorWrapper(r.Path, r.Handler))
	}

	return server
}

func (s *Server) ListenAndServe(ctx context.Context, handlers ...ShutdownHandler) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = s.Server.ListenAndServe()
		interrupt <- syscall.SIGTERM
	}()

	<-interrupt
	if err := s.Shutdown(ctx); err != nil {
		return err
	}

	for _, handler := range handlers {
		err := handler(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Route(r Route) {
	s.routes = append(s.routes, r)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, req)
}
