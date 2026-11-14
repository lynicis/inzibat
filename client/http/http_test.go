package http

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/lynicis/inzibat/internal/testutil"
)

const (
	TestReqPath        = "/items"
	TestReqUri         = "http://localhost"
	TestReqHeaderKey   = "X-Test-Key"
	TestReqHeaderValue = "INZIBAT"
)

var (
	TestReqBody  = []byte(`{"inzibat":"awesome"}`)
	TestRespBody = []byte(`{"status": "ok"}`)
)

func TestClient_Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		var xTestKeyHeader []string
		mockServer.Get(TestReqPath, func(ctx *fiber.Ctx) error {
			xTestKeyHeader = ctx.GetReqHeaders()[TestReqHeaderKey]
			return ctx.Status(fiber.StatusOK).Send(TestReqBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodGet, uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, nil)

		assert.NoError(t, err)
		assert.Equal(t, []string{
			TestReqHeaderValue,
		}, xTestKeyHeader)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestReqBody,
		}, response)
	})

	t.Run("when upstream returns 404", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodGet, uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, fiber.StatusNotFound, response.Status)
	})
}

func TestClient_Post(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		var requestBodyBytes []byte
		mockServer.Post(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodPost, uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when upstream returns 404", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodPost, uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, fiber.StatusNotFound, response.Status)
	})
}

func TestClient_Put(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		var requestBodyBytes []byte
		mockServer.Put(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodPut, uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when upstream returns 404", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodPut, uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, fiber.StatusNotFound, response.Status)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		var requestBodyBytes []byte
		mockServer.Delete(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodDelete, url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when upstream returns 404", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodDelete, url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, fiber.StatusNotFound, response.Status)
	})
}

func TestClient_Patch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		var requestBodyBytes []byte
		mockServer.Patch(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodPatch, url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when upstream returns 404", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		response, err := httpClient.Do(http.MethodPatch, url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, fiber.StatusNotFound, response.Status)
	})
}

func TestClient_SetRetryConfig(t *testing.T) {
	t.Run("happy path - sets retry config", func(t *testing.T) {
		httpClient := NewHttpClient()
		newConfig := RetryConfig{
			MaxRetries:        5,
			InitialBackoff:    200 * time.Millisecond,
			MaxBackoff:        5 * time.Second,
			BackoffMultiplier: 3.0,
		}

		httpClient.SetRetryConfig(newConfig)

		assert.Equal(t, newConfig.MaxRetries, httpClient.retryConfig.MaxRetries)
		assert.Equal(t, newConfig.InitialBackoff, httpClient.retryConfig.InitialBackoff)
		assert.Equal(t, newConfig.MaxBackoff, httpClient.retryConfig.MaxBackoff)
		assert.Equal(t, newConfig.BackoffMultiplier, httpClient.retryConfig.BackoffMultiplier)
	})
}

func TestDefaultRetryConfig(t *testing.T) {
	t.Run("happy path - returns default config", func(t *testing.T) {
		config := DefaultRetryConfig()

		assert.Equal(t, 3, config.MaxRetries)
		assert.Equal(t, 100*time.Millisecond, config.InitialBackoff)
		assert.Equal(t, 2*time.Second, config.MaxBackoff)
		assert.Equal(t, 2.0, config.BackoffMultiplier)
	})
}

func TestClient_calculateBackoff(t *testing.T) {
	t.Run("stays within the jittered range", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        2 * time.Second,
			BackoffMultiplier: 2.0,
		})

		for attempt, expected := range map[int]time.Duration{
			0: 100 * time.Millisecond,
			1: 200 * time.Millisecond,
			2: 400 * time.Millisecond,
		} {
			backoff := httpClient.calculateBackoff(attempt)

			assert.InDelta(t, float64(expected), float64(backoff), float64(expected)/5)
			assert.LessOrEqual(t, backoff, httpClient.retryConfig.MaxBackoff)
		}
	})

	t.Run("happy path - backoff capped at MaxBackoff", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        500 * time.Millisecond,
			BackoffMultiplier: 10.0,
		})

		backoff := httpClient.calculateBackoff(2)

		assert.Equal(t, 500*time.Millisecond, backoff)
	})

	t.Run("jitter varies between calls", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        10 * time.Second,
			BackoffMultiplier: 2.0,
		})

		samples := make(map[time.Duration]struct{})
		for range 100 {
			samples[httpClient.calculateBackoff(3)] = struct{}{}
		}

		assert.Greater(t, len(samples), 1)
	})
}

func TestIsRetryableError(t *testing.T) {
	t.Run("happy path - returns true for error", func(t *testing.T) {
		err := assert.AnError

		result := isRetryableError(err, 0)

		assert.True(t, result)
	})

	t.Run("happy path - returns true for 5xx status codes", func(t *testing.T) {
		assert.True(t, isRetryableError(nil, http.StatusInternalServerError))
		assert.True(t, isRetryableError(nil, http.StatusBadGateway))
		assert.True(t, isRetryableError(nil, http.StatusServiceUnavailable))
		assert.True(t, isRetryableError(nil, 599))
	})

	t.Run("happy path - returns false for 2xx status codes", func(t *testing.T) {
		assert.False(t, isRetryableError(nil, http.StatusOK))
		assert.False(t, isRetryableError(nil, http.StatusCreated))
		assert.False(t, isRetryableError(nil, http.StatusNoContent))
	})

	t.Run("happy path - returns false for 4xx status codes", func(t *testing.T) {
		assert.False(t, isRetryableError(nil, http.StatusBadRequest))
		assert.False(t, isRetryableError(nil, http.StatusNotFound))
		assert.False(t, isRetryableError(nil, http.StatusUnauthorized))
	})

	t.Run("happy path - returns false for status codes >= 600", func(t *testing.T) {
		assert.False(t, isRetryableError(nil, 600))
		assert.False(t, isRetryableError(nil, 700))
	})
}

func TestClient_shouldRetry(t *testing.T) {
	t.Run("happy path - should retry when error and attempt < max", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries: 3,
		})

		assert.True(t, httpClient.shouldRetry(http.MethodGet, assert.AnError, 0, 0))
		assert.True(t, httpClient.shouldRetry(http.MethodGet, assert.AnError, 0, 1))
		assert.True(t, httpClient.shouldRetry(http.MethodGet, assert.AnError, 0, 2))
		assert.False(t, httpClient.shouldRetry(http.MethodGet, assert.AnError, 0, 3))
	})

	t.Run("happy path - should retry for 5xx status codes", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries: 3,
		})

		assert.True(t, httpClient.shouldRetry(http.MethodGet, nil, http.StatusInternalServerError, 0))
		assert.True(t, httpClient.shouldRetry(http.MethodPut, nil, http.StatusBadGateway, 1))
		assert.False(t, httpClient.shouldRetry(http.MethodDelete, nil, http.StatusInternalServerError, 3))
	})

	t.Run("happy path - should not retry for 2xx status codes", func(t *testing.T) {
		httpClient := NewHttpClient()

		assert.False(t, httpClient.shouldRetry(http.MethodGet, nil, http.StatusOK, 0))
		assert.False(t, httpClient.shouldRetry(http.MethodGet, nil, http.StatusCreated, 1))
	})

	t.Run("non-idempotent methods never retry", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries: 3,
		})

		assert.False(t, httpClient.shouldRetry(http.MethodPost, assert.AnError, 0, 0))
		assert.False(t, httpClient.shouldRetry(http.MethodPost, nil, http.StatusInternalServerError, 0))
		assert.False(t, httpClient.shouldRetry(http.MethodPatch, assert.AnError, 0, 0))
		assert.False(t, httpClient.shouldRetry(http.MethodPatch, nil, http.StatusBadGateway, 0))
	})
}

func TestClient_buildRequest(t *testing.T) {
	t.Run("happy path - builds request with headers and body", func(t *testing.T) {
		httpClient := NewHttpClient()
		uri := "http://localhost:8080/test"
		method := http.MethodPost
		headers := http.Header{
			"Content-Type": {"application/json"},
			"X-Custom":     {"value"},
		}
		body := []byte(`{"test": "data"}`)

		req := httpClient.buildRequest(uri, method, headers, body)

		assert.NotNil(t, req)
		assert.Equal(t, uri, string(req.RequestURI()))
		assert.Equal(t, method, string(req.Header.Method()))
		assert.Equal(t, body, req.Body())
	})

	t.Run("happy path - builds request without body", func(t *testing.T) {
		httpClient := NewHttpClient()
		uri := "http://localhost:8080/test"
		method := http.MethodGet
		headers := http.Header{}

		req := httpClient.buildRequest(uri, method, headers, nil)

		assert.NotNil(t, req)
		assert.Equal(t, uri, string(req.RequestURI()))
		assert.Equal(t, method, string(req.Header.Method()))
		assert.Empty(t, req.Body())
	})

	t.Run("sends every value of a multi-value header", func(t *testing.T) {
		httpClient := NewHttpClient()
		headers := http.Header{"X-Multi": {"first", "second"}}

		req := httpClient.buildRequest("http://localhost:8080/test", http.MethodGet, headers, nil)

		assert.Equal(t, [][]byte{[]byte("first"), []byte("second")}, req.Header.PeekAll("X-Multi"))
	})

	t.Run("drops headers owned by the transport", func(t *testing.T) {
		httpClient := NewHttpClient()
		headers := http.Header{
			"Host":              {"example.com"},
			"Content-Length":    {"10"},
			"Connection":        {"keep-alive"},
			"Transfer-Encoding": {"chunked"},
			"Keep-Alive":        {"timeout=5"},
			"X-Kept":            {"value"},
		}

		req := httpClient.buildRequest("http://localhost:8080/test", http.MethodGet, headers, nil)

		assert.NotEqual(t, []byte("example.com"), req.Header.Peek("Host"))
		assert.NotEqual(t, []byte("10"), req.Header.Peek("Content-Length"))
		assert.Empty(t, req.Header.PeekAll("Connection"))
		assert.Empty(t, req.Header.PeekAll("Transfer-Encoding"))
		assert.Empty(t, req.Header.PeekAll("Keep-Alive"))
		assert.Equal(t, []byte("value"), req.Header.Peek("X-Kept"))
	})
}

func TestClient_executeRequest(t *testing.T) {
	t.Run("happy path - executes request successfully", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		mockServer.Get("/test", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/test")
		req := httpClient.buildRequest(uri, http.MethodGet, http.Header{}, nil)

		resp, err := httpClient.executeRequest(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode())
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	})

	t.Run("error path - connection error releases resources", func(t *testing.T) {
		httpClient := NewHttpClient()
		uri := "http://localhost:99999/nonexistent"
		req := httpClient.buildRequest(uri, http.MethodGet, http.Header{}, nil)

		resp, err := httpClient.executeRequest(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestClient_handleResponse(t *testing.T) {
	t.Run("happy path - handles successful response", func(t *testing.T) {
		httpClient := NewHttpClient()
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		resp.SetStatusCode(fiber.StatusOK)
		resp.SetBody(TestRespBody)

		result, statusCode, err := httpClient.handleResponse(resp, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, fiber.StatusOK, result.Status)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, TestRespBody, result.Body)
	})

	t.Run("happy path - handles 3xx redirect response", func(t *testing.T) {
		httpClient := NewHttpClient()
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		resp.SetStatusCode(fiber.StatusMovedPermanently)

		result, statusCode, err := httpClient.handleResponse(resp, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, fiber.StatusMovedPermanently, result.Status)
		assert.Equal(t, fiber.StatusMovedPermanently, statusCode)
	})

	t.Run("happy path - handles 4xx client error response", func(t *testing.T) {
		httpClient := NewHttpClient()
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		resp.SetStatusCode(fiber.StatusBadRequest)

		result, statusCode, err := httpClient.handleResponse(resp, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, fiber.StatusBadRequest, result.Status)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
	})

	t.Run("error path - handles 5xx server error response", func(t *testing.T) {
		httpClient := NewHttpClient()
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		resp.SetStatusCode(fiber.StatusInternalServerError)

		result, statusCode, err := httpClient.handleResponse(resp, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Contains(t, err.Error(), "response failed")
	})
}

func TestClient_makeRequest(t *testing.T) {
	t.Run("happy path - successful request on first attempt", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		})

		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		mockServer.Get("/test", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/test")
		response, err := httpClient.makeRequest(uri, http.MethodGet, http.Header{}, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, fiber.StatusOK, response.Status)
	})

	t.Run("error path - retries on 5xx error and eventually fails", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        1,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        50 * time.Millisecond,
			BackoffMultiplier: 2.0,
		})

		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		mockServer.Get("/test", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusInternalServerError).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/test")
		response, err := httpClient.makeRequest(uri, http.MethodGet, http.Header{}, nil)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error path - retries on connection error and eventually fails", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        1,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        50 * time.Millisecond,
			BackoffMultiplier: 2.0,
		})

		uri := "http://localhost:99999/nonexistent"
		response, err := httpClient.makeRequest(uri, http.MethodGet, http.Header{}, nil)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error path - fails after all retries exhausted", func(t *testing.T) {
		freePort, err := testutil.GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    10 * time.Millisecond,
			MaxBackoff:        50 * time.Millisecond,
			BackoffMultiplier: 2.0,
		})

		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		mockServer.Get("/test", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusInternalServerError).Send(TestRespBody)
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/test")
		response, err := httpClient.makeRequest(uri, http.MethodGet, http.Header{}, nil)

		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestClient_makeRequest_SendsRepeatedHeaders(t *testing.T) {
	freePort, err := testutil.GetFreePort()
	require.NoError(t, err)

	mockServer := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	captured := make(chan []string, 1)
	mockServer.Get("/multi", func(ctx *fiber.Ctx) error {
		captured <- ctx.GetReqHeaders()["X-Multi"]

		return ctx.SendStatus(fiber.StatusOK)
	})

	go mockServer.Listen(fmt.Sprintf(":%d", freePort))
	defer mockServer.Shutdown()
	time.Sleep(1 * time.Second)

	uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/multi")
	_, err = NewHttpClient().makeRequest(uri, http.MethodGet, http.Header{
		"X-Multi": {"first", "second"},
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, []string{"first", "second"}, <-captured)
}

func TestClient_makeRequest_RetryMethodMatrix(t *testing.T) {
	freePort, err := testutil.GetFreePort()
	require.NoError(t, err)

	var getAttempts, postAttempts atomic.Int32

	mockServer := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	mockServer.Get("/fail", func(ctx *fiber.Ctx) error {
		getAttempts.Add(1)
		return ctx.Status(fiber.StatusInternalServerError).SendString("boom")
	})
	mockServer.Post("/fail", func(ctx *fiber.Ctx) error {
		postAttempts.Add(1)
		return ctx.Status(fiber.StatusInternalServerError).SendString("boom")
	})

	go mockServer.Listen(fmt.Sprintf(":%d", freePort))
	defer mockServer.Shutdown()
	time.Sleep(1 * time.Second)

	uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/fail")

	httpClient := NewHttpClient()
	httpClient.SetRetryConfig(RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    time.Millisecond,
		MaxBackoff:        5 * time.Millisecond,
		BackoffMultiplier: 2.0,
	})

	_, err = httpClient.makeRequest(uri, http.MethodGet, http.Header{}, nil)
	assert.Error(t, err)
	assert.Equal(t, int32(3), getAttempts.Load())

	_, err = httpClient.makeRequest(uri, http.MethodPost, http.Header{}, []byte("{}"))
	assert.Error(t, err)
	assert.Equal(t, int32(1), postAttempts.Load())
}

func TestClient_makeRequest_RetriesOn5xxThenSucceeds(t *testing.T) {
	freePort, err := testutil.GetFreePort()
	require.NoError(t, err)

	httpClient := NewHttpClient()
	httpClient.SetRetryConfig(RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    time.Millisecond,
		MaxBackoff:        5 * time.Millisecond,
		BackoffMultiplier: 2.0,
	})

	mockServer := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	var attempts atomic.Int32
	mockServer.Get("/flaky", func(ctx *fiber.Ctx) error {
		if attempts.Add(1) == 1 {
			return ctx.Status(fiber.StatusInternalServerError).SendString("boom")
		}

		return ctx.Status(fiber.StatusTeapot).SendString("recovered")
	})

	go mockServer.Listen(fmt.Sprintf(":%d", freePort))
	defer mockServer.Shutdown()
	time.Sleep(1 * time.Second)

	uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, "/flaky")
	response, err := httpClient.makeRequest(uri, http.MethodGet, http.Header{}, nil)

	require.NoError(t, err)
	assert.Equal(t, int32(2), attempts.Load())
	assert.Equal(t, fiber.StatusTeapot, response.Status)
	assert.Equal(t, []byte("recovered"), response.Body)
}

func TestIsUnforwardableHeader(t *testing.T) {
	for _, name := range []string{
		"Host", "Content-Length", "Connection", "Transfer-Encoding", "Keep-Alive",
		"Upgrade", "Proxy-Connection", "Te", "Trailer", "Proxy-Authorization",
		"host", "transfer-encoding",
	} {
		assert.True(t, IsUnforwardableHeader(name), "header: %s", name)
	}

	for _, name := range []string{"X-Test", "Content-Type", "Authorization"} {
		assert.False(t, IsUnforwardableHeader(name), "header: %s", name)
	}
}
