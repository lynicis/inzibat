package client_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/port"
	"github.com/Lynicis/inzibat/server"

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
	c := client.NewClient()
	assert.Implements(t, (*client.Client)(nil), c)
}

func TestClient_Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		var xTestKeyHeader string
		mockServer.GetFiberInstance().Get(TestReqPath, func(ctx *fiber.Ctx) error {
			xTestKeyHeader = ctx.GetReqHeaders()[TestReqHeaderKey]
			return ctx.Status(fiber.StatusOK).Send(TestReqBody)
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		response, err := c.Get(uri, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		})

		assert.NoError(t, err)
		assert.Equal(t, TestReqHeaderValue, xTestKeyHeader)
		assert.Equal(t, &client.HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestReqBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Get(uri, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		})

		assert.Error(t, err)
	})
}

func TestClient_Post(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
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
		response, err := c.Post(uri, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &client.HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Post(uri, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Put(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
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
		response, err := c.Put(uri, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &client.HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		uri := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Put(uri, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
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
		response, err := c.Delete(url, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &client.HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Delete(url, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_Patch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
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
		response, err := c.Patch(url, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.NoError(t, err)
		assert.Equal(t, TestReqBody, requestBodyBytes)
		assert.Equal(t, &client.HttpResponse{
			Status: fiber.StatusOK,
			Body:   TestRespBody,
		}, response)
	})

	t.Run("when client return error", func(t *testing.T) {
		freePort, err := port.GetFreePort()
		require.NoError(t, err)

		c := client.NewClient()
		mockServer := server.NewServer(&config.Config{
			ServerPort: freePort,
		})

		go mockServer.Start()
		defer mockServer.Shutdown()
		time.Sleep(1 * time.Second)

		url := fmt.Sprintf("%s:%s%s", TestReqUri, freePort, TestReqPath)
		_, err = c.Patch(url, client.HttpHeader{
			TestReqHeaderKey: TestReqHeaderValue,
		}, TestReqBody)

		assert.Error(t, err)
	})
}

func TestClient_GetCloneOfStruct(t *testing.T) {
	newClient := client.NewClient()
	cloneOfStruct := newClient.GetCloneOfStruct()

	assert.NotSame(t, newClient, cloneOfStruct)
	assert.Implements(t, (*client.Client)(nil), cloneOfStruct)
}
