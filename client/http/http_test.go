package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		freePort, err := GetFreePort()
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
		response, err := httpClient.Get(uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{
			TestReqHeaderValue,
		}, xTestKeyHeader)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestReqBody,
		}, response)
	})

	t.Run("when HttpClient return error", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Get(uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		})

		assert.Error(t, err)
	})
}

func TestClient_Post(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
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
		response, err := httpClient.Post(uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when HttpClient return error", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Post(uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Put(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
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
		response, err := httpClient.Put(uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when HttpClient return error", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Put(uri, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
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
		response, err := httpClient.Delete(url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when HttpClient return error", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Delete(url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Patch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
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
		response, err := httpClient.Patch(url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &Response{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when HttpClient return error", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := NewHttpClient()
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Patch(url, http.Header{
			TestReqHeaderKey: {TestReqHeaderValue},
		}, TestReqBody)

		assert.Error(t, err)
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
	t.Run("happy path - calculates backoff correctly", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries:        3,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        2 * time.Second,
			BackoffMultiplier: 2.0,
		})

		backoff1 := httpClient.calculateBackoff(0)
		assert.Equal(t, 100*time.Millisecond, backoff1)

		backoff2 := httpClient.calculateBackoff(1)
		assert.Equal(t, 200*time.Millisecond, backoff2)

		backoff3 := httpClient.calculateBackoff(2)
		assert.Equal(t, 400*time.Millisecond, backoff3)
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
}

func TestPow(t *testing.T) {
	t.Run("happy path - calculates power correctly", func(t *testing.T) {
		assert.Equal(t, 1.0, pow(2.0, 0))
		assert.Equal(t, 2.0, pow(2.0, 1))
		assert.Equal(t, 4.0, pow(2.0, 2))
		assert.Equal(t, 8.0, pow(2.0, 3))
		assert.Equal(t, 9.0, pow(3.0, 2))
	})

	t.Run("happy path - handles zero base", func(t *testing.T) {
		assert.Equal(t, 1.0, pow(0.0, 0))
		assert.Equal(t, 0.0, pow(0.0, 1))
	})

	t.Run("happy path - handles zero exponent", func(t *testing.T) {
		assert.Equal(t, 1.0, pow(5.0, 0))
		assert.Equal(t, 1.0, pow(10.0, 0))
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

		assert.True(t, httpClient.shouldRetry(assert.AnError, 0, 0))
		assert.True(t, httpClient.shouldRetry(assert.AnError, 0, 1))
		assert.True(t, httpClient.shouldRetry(assert.AnError, 0, 2))
		assert.False(t, httpClient.shouldRetry(assert.AnError, 0, 3))
	})

	t.Run("happy path - should retry for 5xx status codes", func(t *testing.T) {
		httpClient := NewHttpClient()
		httpClient.SetRetryConfig(RetryConfig{
			MaxRetries: 3,
		})

		assert.True(t, httpClient.shouldRetry(nil, http.StatusInternalServerError, 0))
		assert.True(t, httpClient.shouldRetry(nil, http.StatusBadGateway, 1))
		assert.False(t, httpClient.shouldRetry(nil, http.StatusInternalServerError, 3))
	})

	t.Run("happy path - should not retry for 2xx status codes", func(t *testing.T) {
		httpClient := NewHttpClient()

		assert.False(t, httpClient.shouldRetry(nil, http.StatusOK, 0))
		assert.False(t, httpClient.shouldRetry(nil, http.StatusCreated, 1))
	})
}

func TestGetFreePort(t *testing.T) {
	t.Run("happy path - returns a free port", func(t *testing.T) {
		port, err := GetFreePort()

		assert.NoError(t, err)
		assert.Greater(t, port, 0)
		assert.LessOrEqual(t, port, 65535)
	})

	t.Run("happy path - returns different ports on multiple calls", func(t *testing.T) {
		port1, err1 := GetFreePort()
		port2, err2 := GetFreePort()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Greater(t, port1, 0)
		assert.Greater(t, port2, 0)
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
}
