package middlewares

import (
	"net/http"
)

func Middlewares(main http.Handler, middlewares ...func(handler http.Handler) http.Handler) http.Handler {
	handler := main
	for i := range middlewares {
		handler = middlewares[len(middlewares)-1-i](handler)
	}

	return handler
}
