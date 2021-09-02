package server

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/diegodesousas/httpserver/pkg/httprouter"

	"github.com/diegodesousas/httpserver/pkg/chi"

	"github.com/stretchr/testify/mock"

	mocks "github.com/diegodesousas/httpserver/mocks/server"

	"github.com/stretchr/testify/assert"
)

type TestRoundTripper struct {
	Retry int
}

func (t *TestRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil && t.Retry > 0 {
		t.Retry--
		time.Sleep(time.Millisecond)
		return t.RoundTrip(request)
	}

	return response, err
}

func TestNew(t *testing.T) {
	assertions := assert.New(t)

	t.Run("should create a server successfully", func(t *testing.T) {
		route := Route{
			Path:    "/test",
			Method:  http.MethodGet,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, res *http.Request) {}),
		}

		mockRouter := new(mocks.Router)
		mockRouter.On("Resource", mock.Anything, mock.Anything, mock.Anything)

		server := New(
			WithPort(8000),
			WithRoutes(route),
			WithRouter(mockRouter),
		)

		assertions.NotNil(server)
	})

	t.Run("should shutdown server success successfully", func(t *testing.T) {
		route := Route{
			Path:    "/test",
			Method:  http.MethodGet,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, res *http.Request) {}),
		}

		mockRouter := new(mocks.Router)
		mockRouter.On("Resource", mock.Anything, mock.Anything, mock.Anything)

		server := New(
			WithPort(8001),
			WithRoutes(route),
			WithRouter(mockRouter),
		)

		assertions.NotNil(server)

		go func() {
			time.Sleep(time.Millisecond)
			err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			assertions.Nil(err)
		}()

		shutdownCount := 0
		shutdownFunc := func(ctx context.Context) error {
			shutdownCount++
			return nil
		}

		err := server.ListenAndServe(context.Background(), shutdownFunc)
		assertions.Nil(err)
		assertions.Equal(1, shutdownCount)
	})

	t.Run("should shutdown server with shutdown handler error", func(t *testing.T) {
		route := Route{
			Path:    "/test",
			Method:  http.MethodGet,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, res *http.Request) {}),
		}

		mockRouter := new(mocks.Router)
		mockRouter.On("Resource", mock.Anything, mock.Anything, mock.Anything)

		server := New(
			WithPort(8001),
			WithRoutes(route),
			WithRouter(mockRouter),
		)

		assertions.NotNil(server)

		go func() {
			time.Sleep(time.Millisecond)
			err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			assertions.Nil(err)
		}()

		expectedErr := errors.New("expected error")
		shutdownFunc := func(ctx context.Context) error {
			return expectedErr
		}

		err := server.ListenAndServe(context.Background(), shutdownFunc)
		assertions.NotNil(err)
		assertions.Equal(expectedErr, err)
	})

	t.Run("should create a server and listen successfully", func(t *testing.T) {
		route := Route{
			Path:    "/test",
			Method:  http.MethodGet,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, res *http.Request) {}),
		}

		server := New(
			WithPort(8001),
			WithRoutes(route),
			WithRouter(chi.New()),
		)

		assertions.NotNil(server)

		go func() {
			err := server.ListenAndServe(context.Background())
			assertions.Nil(err)
		}()

		client := http.DefaultClient
		client.Transport = &TestRoundTripper{Retry: 5}

		res, err := client.Get("http://localhost:8001/test")
		assertions.Nil(err)
		assertions.Equal(http.StatusOK, res.StatusCode)

		err = server.Shutdown(context.Background())
		assertions.Nil(err)
	})

	t.Run("should create a server and listen successfully with middlewares", func(t *testing.T) {
		route := Route{
			Path:    "/test",
			Method:  http.MethodGet,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, res *http.Request) {}),
		}

		orderedMessage := ""
		m1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				orderedMessage = "first"
				next.ServeHTTP(w, req)
			})
		}

		m2 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				orderedMessage = orderedMessage + " second"
				next.ServeHTTP(w, req)
			})
		}

		server := New(
			WithPort(8001),
			WithRoutes(route),
			WithRouter(httprouter.New()),
			WithMiddlewares(m1, m2),
		)

		assertions.NotNil(server)

		go func() {
			err := server.ListenAndServe(context.Background())
			assertions.Nil(err)
		}()

		client := http.DefaultClient
		client.Transport = &TestRoundTripper{Retry: 5}

		res, err := client.Get("http://localhost:8001/test")
		assertions.Nil(err)
		assertions.Equal(http.StatusOK, res.StatusCode)
		assertions.Equal("first second", orderedMessage)

		err = server.Shutdown(context.Background())
		assertions.Nil(err)
	})
}
