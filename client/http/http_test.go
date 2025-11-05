package http

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
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

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
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
		response, err := httpClient.Get(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
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

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Get(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		})

		assert.Error(t, err)
	})
}

func TestClient_Post(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
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
		response, err := httpClient.Post(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
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

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Post(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Put(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
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
		response, err := httpClient.Put(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
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

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Put(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
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
		response, err := httpClient.Delete(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
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

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Delete(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Patch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := GetFreePort()
		require.NoError(t, err)

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
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
		response, err := httpClient.Patch(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
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

		httpClient := &HttpClient{
			client: &fasthttp.Client{},
		}
		mockServer := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go mockServer.Listen(fmt.Sprintf(":%d", freePort))
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%d%s", TestReqUri, freePort, TestReqPath)
		_, err = httpClient.Patch(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer tcpListener.Close()

	return tcpListener.Addr().(*net.TCPAddr).Port, nil
}
