package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpPkg "github.com/lynicis/inzibat/client/http"
	"github.com/lynicis/inzibat/config"
)

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
	})
}
