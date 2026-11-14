package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpPkg "github.com/lynicis/inzibat/client/http"
	"github.com/lynicis/inzibat/config"
)

func TestClientHandler_allowRequest(t *testing.T) {
	t.Run("happy path - no circuit breaker", func(t *testing.T) {
		clientHandler := &ClientHandler{}

		allowed, err := clientHandler.allowRequest(false, "")

		assert.NoError(t, err)
		assert.True(t, allowed)
	})
}

func TestClientHandler_recordSuccess(t *testing.T) {
	t.Run("happy path - no circuit breaker", func(t *testing.T) {
		clientHandler := &ClientHandler{}

		err := clientHandler.recordSuccess(false, "")

		assert.NoError(t, err)
	})

	t.Run("happy path - with circuit breaker", func(t *testing.T) {
		circuitBreakerStore := NewCircuitBreakerStore()

		routeKey := "GET /proxy -> GET http://example.com/"
		require.NoError(t, circuitBreakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:          config.BoolPointer(true),
			FailureThreshold: 1,
			MinimumRequests:  1,
		}))

		clientHandler := &ClientHandler{
			CircuitBreakerStore: circuitBreakerStore,
		}

		err := clientHandler.recordSuccess(true, routeKey)

		assert.NoError(t, err)
	})

	t.Run("error path - circuit breaker store returns error", func(t *testing.T) {
		circuitBreakerStore := NewCircuitBreakerStore()

		clientHandler := &ClientHandler{
			CircuitBreakerStore: circuitBreakerStore,
		}

		err := clientHandler.recordSuccess(true, "non-existent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update circuit breaker success")
	})
}

func TestClientHandler_recordFailure(t *testing.T) {
	t.Run("happy path - no circuit breaker", func(t *testing.T) {
		clientHandler := &ClientHandler{}

		err := clientHandler.recordFailure(false, "")

		assert.NoError(t, err)
	})

	t.Run("error path - circuit breaker store returns error", func(t *testing.T) {
		circuitBreakerStore := NewCircuitBreakerStore()

		clientHandler := &ClientHandler{
			CircuitBreakerStore: circuitBreakerStore,
		}

		err := clientHandler.recordFailure(true, "non-existent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update circuit breaker failure")
	})
}

func TestClientHandler_CreateHandler(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Run("GET method", func(t *testing.T) {
			targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message":"success"}`))
			}))
			defer targetServer.Close()

			httpClient := httpPkg.NewHttpClient()
			clientHandler := &ClientHandler{
				Client: httpClient,
				RouteConfig: &[]config.Route{
					{
						Method: http.MethodGet,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method: http.MethodGet,
							Headers: http.Header{
								"X-Custom-Header": {"test-value"},
							},
							Host:                   targetServer.URL,
							Path:                   "/",
							PassWithRequestBody:    false,
							PassWithRequestHeaders: false,
							InErrorReturn500:       false,
						},
					},
				},
			}

			handler := clientHandler.CreateHandler(0)
			fiberApp := fiber.New()
			fiberApp.Get("/proxy", handler)

			request := httptest.NewRequest(http.MethodGet, "/proxy", nil)
			response, err := fiberApp.Test(request)

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Contains(t, string(responseBody), "success")
		})

		t.Run("POST method", func(t *testing.T) {
			targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"id":123}`))
			}))
			defer targetServer.Close()

			httpClient := httpPkg.NewHttpClient()
			clientHandler := &ClientHandler{
				Client: httpClient,
				RouteConfig: &[]config.Route{
					{
						Method: http.MethodPost,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method: http.MethodPost,
							Headers: http.Header{
								"Content-Type": {"application/json"},
							},
							Body: config.HttpBody{
								"name": "test",
							},
							Host:                   targetServer.URL,
							Path:                   "/",
							PassWithRequestBody:    false,
							PassWithRequestHeaders: false,
							InErrorReturn500:       false,
						},
					},
				},
			}

			handler := clientHandler.CreateHandler(0)
			fiberApp := fiber.New()
			fiberApp.Post("/proxy", handler)

			request := httptest.NewRequest(http.MethodPost, "/proxy", nil)
			response, err := fiberApp.Test(request)

			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, response.StatusCode)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Contains(t, string(responseBody), "123")
		})
	})

	t.Run("error scenarios", func(t *testing.T) {
		t.Run("invalid URL parsing", func(t *testing.T) {
			httpClient := httpPkg.NewHttpClient()
			clientHandler := &ClientHandler{
				Client: httpClient,
				RouteConfig: &[]config.Route{
					{
						Method: http.MethodGet,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method: http.MethodGet,
							Host:   "://invalid-url",
							Path:   "/test",
						},
					},
				},
			}

			handler := clientHandler.CreateHandler(0)
			fiberApp := fiber.New()
			fiberApp.Get("/proxy", handler)

			request := httptest.NewRequest(http.MethodGet, "/proxy", nil)
			response, err := fiberApp.Test(request)

			require.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		})

		t.Run("HTTP client error with InErrorReturn500 false", func(t *testing.T) {
			httpClient := httpPkg.NewHttpClient()
			clientHandler := &ClientHandler{
				Client: httpClient,
				RouteConfig: &[]config.Route{
					{
						Method: http.MethodGet,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method:                 http.MethodGet,
							Host:                   "http://127.0.0.1:99999",
							Path:                   "/test",
							PassWithRequestBody:    false,
							PassWithRequestHeaders: false,
							InErrorReturn500:       false,
						},
					},
				},
			}

			handler := clientHandler.CreateHandler(0)
			fiberApp := fiber.New()
			fiberApp.Get("/proxy", handler)

			request := httptest.NewRequest(http.MethodGet, "/proxy", nil)
			response, err := fiberApp.Test(request)

			require.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, string(responseBody))
		})

		t.Run("HTTP client error with InErrorReturn500 true", func(t *testing.T) {
			httpClient := httpPkg.NewHttpClient()
			clientHandler := &ClientHandler{
				Client: httpClient,
				RouteConfig: &[]config.Route{
					{
						Method: http.MethodGet,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method:                 http.MethodGet,
							Host:                   "http://127.0.0.1:99999",
							Path:                   "/test",
							PassWithRequestBody:    false,
							PassWithRequestHeaders: false,
							InErrorReturn500:       true,
						},
					},
				},
			}

			handler := clientHandler.CreateHandler(0)
			fiberApp := fiber.New()
			fiberApp.Get("/proxy", handler)

			request := httptest.NewRequest(http.MethodGet, "/proxy", nil)
			response, err := fiberApp.Test(request)

			require.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Empty(t, string(responseBody))
		})

		t.Run("circuit breaker open returns 503 without upstream call", func(t *testing.T) {
			var upstreamCallCount int32
			targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&upstreamCallCount, 1)
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer targetServer.Close()

			httpClient := httpPkg.NewHttpClient()
			httpClient.SetRetryConfig(httpPkg.RetryConfig{
				MaxRetries:        0,
				InitialBackoff:    1,
				MaxBackoff:        1,
				BackoffMultiplier: 1,
			})

			clientHandler := &ClientHandler{
				Client: httpClient,
				RouteConfig: &[]config.Route{
					{
						Method: http.MethodGet,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method: http.MethodGet,
							Host:   targetServer.URL,
							Path:   "/",
						},
					},
				},
			}

			circuitBreakerStore := NewCircuitBreakerStore()

			routeKey := "GET /proxy -> GET " + targetServer.URL + "/"
			err := circuitBreakerStore.Seed(routeKey, config.CircuitBreakerConfig{
				Enabled:             config.BoolPointer(true),
				FailureThreshold:    1,
				MinimumRequests:     1,
				OpenTimeoutMs:       60000,
				HalfOpenMaxRequests: 1,
				SuccessThreshold:    1,
			})
			require.NoError(t, err)

			clientHandler.CircuitBreakerStore = circuitBreakerStore
			clientHandler.CircuitBreakerRouteKeys = map[int]string{0: routeKey}

			handler := clientHandler.CreateHandler(0)
			fiberApp := fiber.New()
			fiberApp.Get("/proxy", handler)

			firstRequest := httptest.NewRequest(http.MethodGet, "/proxy", nil)
			firstResponse, err := fiberApp.Test(firstRequest)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, firstResponse.StatusCode)
			assert.Equal(t, int32(1), atomic.LoadInt32(&upstreamCallCount))

			secondRequest := httptest.NewRequest(http.MethodGet, "/proxy", nil)
			secondResponse, err := fiberApp.Test(secondRequest)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusServiceUnavailable, secondResponse.StatusCode)
			assert.Equal(t, int32(1), atomic.LoadInt32(&upstreamCallCount))
		})
	})
}

func TestClientHandler_CreateHandler_NoNullBody(t *testing.T) {
	captureBody := func(t *testing.T) (*httptest.Server, chan string) {
		t.Helper()
		bodies := make(chan string, 1)
		targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			bodies <- string(body)
			w.WriteHeader(http.StatusOK)
		}))
		t.Cleanup(targetServer.Close)

		return targetServer, bodies
	}

	t.Run("GET sends zero bytes", func(t *testing.T) {
		targetServer, bodies := captureBody(t)

		clientHandler := &ClientHandler{
			Client: httpPkg.NewHttpClient(),
			RouteConfig: &[]config.Route{
				{
					Method: http.MethodGet,
					Path:   "/proxy",
					RequestTo: &config.RequestTo{
						Method: http.MethodGet,
						Host:   targetServer.URL,
						Path:   "/",
					},
				},
			},
		}

		fiberApp := fiber.New()
		fiberApp.Get("/proxy", clientHandler.CreateHandler(0))

		response, err := fiberApp.Test(httptest.NewRequest(http.MethodGet, "/proxy", nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "", <-bodies)
	})

	t.Run("POST with empty configured body sends zero bytes", func(t *testing.T) {
		targetServer, bodies := captureBody(t)

		clientHandler := &ClientHandler{
			Client: httpPkg.NewHttpClient(),
			RouteConfig: &[]config.Route{
				{
					Method: http.MethodPost,
					Path:   "/proxy",
					RequestTo: &config.RequestTo{
						Method: http.MethodPost,
						Host:   targetServer.URL,
						Path:   "/",
					},
				},
			},
		}

		fiberApp := fiber.New()
		fiberApp.Post("/proxy", clientHandler.CreateHandler(0))

		response, err := fiberApp.Test(httptest.NewRequest(http.MethodPost, "/proxy", nil))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "", <-bodies)
	})
}

func TestClientHandler_CreateHandler_PassThrough(t *testing.T) {
	t.Run("passWithRequestHeaders forwards incoming headers", func(t *testing.T) {
		headers := make(chan http.Header, 1)
		targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers <- r.Header.Clone()
			w.WriteHeader(http.StatusOK)
		}))
		defer targetServer.Close()

		clientHandler := &ClientHandler{
			Client: httpPkg.NewHttpClient(),
			RouteConfig: &[]config.Route{
				{
					Method: http.MethodGet,
					Path:   "/proxy",
					RequestTo: &config.RequestTo{
						Method:                 http.MethodGet,
						Headers:                http.Header{"X-Configured": {"configured"}},
						Host:                   targetServer.URL,
						Path:                   "/",
						PassWithRequestHeaders: true,
					},
				},
			},
		}

		fiberApp := fiber.New()
		fiberApp.Get("/proxy", clientHandler.CreateHandler(0))

		request := httptest.NewRequest(http.MethodGet, "/proxy", nil)
		request.Header.Set("X-Test", "1")
		request.Header.Set("Content-Length", "0")

		response, err := fiberApp.Test(request)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		forwarded := <-headers
		assert.Equal(t, "1", forwarded.Get("X-Test"))
		assert.Equal(t, "configured", forwarded.Get("X-Configured"))
	})

	t.Run("passWithRequestBody forwards the incoming body verbatim", func(t *testing.T) {
		bodies := make(chan string, 1)
		targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			bodies <- string(body)
			w.WriteHeader(http.StatusOK)
		}))
		defer targetServer.Close()

		clientHandler := &ClientHandler{
			Client: httpPkg.NewHttpClient(),
			RouteConfig: &[]config.Route{
				{
					Method: http.MethodPost,
					Path:   "/proxy",
					RequestTo: &config.RequestTo{
						Method:              http.MethodPost,
						Body:                config.HttpBody{"configured": "body"},
						Host:                targetServer.URL,
						Path:                "/",
						PassWithRequestBody: true,
					},
				},
			},
		}

		fiberApp := fiber.New()
		fiberApp.Post("/proxy", clientHandler.CreateHandler(0))

		request := httptest.NewRequest(http.MethodPost, "/proxy", strings.NewReader(`{"a":1}`))
		request.Header.Set("Content-Type", "application/json")

		response, err := fiberApp.Test(request)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, `{"a":1}`, <-bodies)
	})
}
