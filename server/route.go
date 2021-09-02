package server

import "net/http"

type Route struct {
	Path    string
	Method  string
	Handler http.Handler
}

type Router interface {
	http.Handler
	Resource(method string, path string, handler http.Handler)
}
