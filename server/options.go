package server

import (
	"fmt"
	"net/http"

	"github.com/diegodesousas/httpserver/middlewares"
)

type Option func(server *Server)

func WithRouter(router Router) Option {
	return func(server *Server) {
		server.router = router
	}
}

func WithRoutes(routes ...Route) Option {
	return func(server *Server) {
		for _, route := range routes {
			server.Route(route)
		}
	}
}

func WithMiddlewares(list ...func(handler http.Handler) http.Handler) Option {
	return func(server *Server) {
		server.Server.Handler = middlewares.Middlewares(
			server,
			list...,
		)
	}
}

func WithPort(port int) Option {
	return func(server *Server) {
		server.Server.Addr = fmt.Sprintf(":%d", port)
	}
}
