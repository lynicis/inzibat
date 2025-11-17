package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func TestMockRoute_CreateRoute(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Run("GET Method", func(t *testing.T) {
			mockRoute := &EndpointHandler{
				RouteConfig: &[]config.Route{
					{
						Method: fiber.MethodGet,
						Path:   "/route-one",
						FakeResponse: &config.FakeResponse{
							Headers: http.Header{
								"X-Test-Header": {"test-header-value"},
							},
							Body: config.HttpBody{
								"nick":     "lynicis",
								"password": 12345,
							},
							StatusCode: 200,
						},
					},
				},
			}
			handler := mockRoute.CreateHandler(0)

			fiberApp := fiber.New()
			fiberApp.Get("/user", handler)

			request := httptest.NewRequest(fiber.MethodGet, "/user", nil)
			response, err := fiberApp.Test(request)
			require.NoError(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			var body map[string]interface{}
			err = json.Unmarshal(responseBody, &body)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusOK, response.StatusCode)
			assert.Equal(t, "lynicis", body["nick"])
			assert.Equal(t, float64(12345), body["password"])
			assert.Equal(t, "test-header-value", response.Header["X-Test-Header"][0])
		})

		t.Run("Other HTTP Method", func(t *testing.T) {
			mockRoute := &EndpointHandler{
				RouteConfig: &[]config.Route{
					{
						Method: fiber.MethodGet,
						Path:   "/route-one",
						FakeResponse: &config.FakeResponse{
							Headers: http.Header{
								"X-Test-Header": {"test-header-value"},
							},
							Body: config.HttpBody{
								"token": "abcd.abcd.abcd",
							},
							StatusCode: 201,
						},
					},
				},
			}
			handler := mockRoute.CreateHandler(0)

			fiberApp := fiber.New()
			fiberApp.Post("/user", handler)

			requestBody := bytes.NewBufferString(`{"nick":"lynicis","password":"1234"}"`)
			request := httptest.NewRequest(fiber.MethodPost, "/user", requestBody)
			response, err := fiberApp.Test(request)
			require.NoError(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			var body map[string]interface{}
			err = json.Unmarshal(responseBody, &body)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusCreated, response.StatusCode)
			assert.Equal(t, "abcd.abcd.abcd", body["token"])
			assert.Equal(t, "test-header-value", response.Header["X-Test-Header"][0])
		})
	})
}
