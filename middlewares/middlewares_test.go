package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewares(t *testing.T) {
	assertions := assert.New(t)

	t.Run("should run middlewares in the specified order", func(t *testing.T) {
		orderedMessage := ""
		main := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			orderedMessage = orderedMessage + "third"
		})

		middleware1 := func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				orderedMessage = orderedMessage + "first "
				handler.ServeHTTP(w, req)
			})
		}

		middleware2 := func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				orderedMessage = orderedMessage + "second "
				handler.ServeHTTP(w, req)
			})
		}

		handler := Middlewares(main, middleware1, middleware2)

		w := httptest.NewRecorder()
		req := &http.Request{}

		handler.ServeHTTP(w, req)

		assertions.Equal("first second third", orderedMessage)
	})
}
