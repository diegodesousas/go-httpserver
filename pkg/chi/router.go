package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Chi struct {
	*chi.Mux
}

func New() *Chi {
	return &Chi{
		chi.NewRouter(),
	}
}

func (c *Chi) Resource(method string, path string, handler http.Handler) {
	c.Method(method, path, handler)
}
