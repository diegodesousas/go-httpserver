# Basic usage example

```
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/diegodesousas/httpserver/pkg/httprouter"
    "github.com/diegodesousas/httpserver/server"
)

func main() {

    ctx := context.Background()
	testRoute := server.Route{
	    Path:   "/test",
	    Method: http.MethodGet,
	    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		    log.Print("Ok")

		    _, err := w.Write([]byte("Ok"))
			if err != nil {
				return
			}
		}),
	}

	customMidd1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			log.Print("custom middleware 1")

			next.ServeHTTP(writer, request)
		})
	}

	customMidd2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			log.Print("custom middleware 2")

			next.ServeHTTP(writer, request)
		})
	}

	s := server.New(
		server.WithPort(8081),
		server.WithRouter(httprouter.New()),
		server.WithRoutes(testRoute),
		server.WithMiddlewares(customMidd1, customMidd2),
	)

	if err := s.ListenAndServe(ctx); err != nil {
		log.Fatal(err)
	}
}
```