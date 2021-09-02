package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	assertions := assert.New(t)

	t.Run("should return status 500 when panic happens", func(t *testing.T) {
		main := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			panic("unexpected error")
		})

		handler := Middlewares(main, PanicRecoveryMiddleware)

		w := httptest.NewRecorder()
		req := &http.Request{}

		handler.ServeHTTP(w, req)

		assertions.Equal(http.StatusInternalServerError, w.Code)
	})
}
