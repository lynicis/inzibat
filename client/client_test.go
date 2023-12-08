package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/server"
	testUtils "github.com/Lynicis/inzibat/test-utils"

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

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.Implements(t, (*Client)(nil), c)
}

func TestClient_Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		var xTestKeyHeader []string
		mockServer.GetFiberInstance().Get(TestReqPath, func(ctx *fiber.Ctx) error {
			xTestKeyHeader = ctx.GetReqHeaders()[TestReqHeaderKey]
			return ctx.Status(fiber.StatusOK).Send(TestReqBody)
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		response, err := c.Get(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{
			TestReqHeaderValue,
		}, xTestKeyHeader)
		assert.Equal(t, &HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestReqBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Get(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		})

		assert.Error(t, err)
	})
}

func TestClient_Post(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		var requestBodyBytes []byte
		mockServer.GetFiberInstance().Post(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		response, err := c.Post(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Post(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Put(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		var requestBodyBytes []byte
		mockServer.GetFiberInstance().Put(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		response, err := c.Put(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Put(uri, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		var requestBodyBytes []byte
		mockServer.GetFiberInstance().Delete(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})
		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		response, err := c.Delete(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Delete(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Patch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		var requestBodyBytes []byte
		mockServer.GetFiberInstance().Patch(TestReqPath, func(ctx *fiber.Ctx) error {
			requestBodyBytes = ctx.Body()
			return ctx.Status(fiber.StatusOK).Send(TestRespBody)
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		response, err := c.Patch(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := testUtils.GetFreePort()
		require.NoError(t, err)

		c := NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Patch(url, map[string]string{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}
