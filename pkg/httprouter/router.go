package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HttpRouter struct {
	*httprouter.Router
}

func New() *HttpRouter {
	return &HttpRouter{
		httprouter.New(),
	}
}

func (h *HttpRouter) Resource(method string, path string, handler http.Handler) {
	h.Handler(method, path, handler)
}
