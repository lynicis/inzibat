package recorder

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRecorderMiddleware(t *testing.T) {
	t.Run("records request and response", func(t *testing.T) {
		store := NewStore(100)
		app := fiber.New()
		app.Use(NewRecorderMiddleware(store))
		app.Get("/hello", func(c *fiber.Ctx) error {
			return c.Status(200).JSON(fiber.Map{"message": "world"})
		})

		req := httptest.NewRequest("GET", "/hello", nil)
		req.Header.Set("X-Test", "value")

		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 1, store.Len())

		entries := store.List()
		entry := entries[0]
		assert.Equal(t, "GET", entry.Request.Method)
		assert.Equal(t, "/hello", entry.Request.Path)
		assert.Equal(t, 200, entry.Response.StatusCode)
		assert.NotEmpty(t, entry.ID)
		assert.NotZero(t, entry.Timestamp)
		assert.GreaterOrEqual(t, entry.DurationMs, int64(0))
	})

	t.Run("captures request body", func(t *testing.T) {
		store := NewStore(100)
		app := fiber.New()
		app.Use(NewRecorderMiddleware(store))
		app.Post("/data", func(c *fiber.Ctx) error {
			return c.Status(201).SendString("created")
		})

		body := `{"key":"value"}`
		req := httptest.NewRequest("POST", "/data", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		entries := store.List()
		require.Len(t, entries, 1)
		assert.Equal(t, `{"key":"value"}`, string(entries[0].Request.Body))
	})

	t.Run("captures response body", func(t *testing.T) {
		store := NewStore(100)
		app := fiber.New()
		app.Use(NewRecorderMiddleware(store))
		app.Get("/json", func(c *fiber.Ctx) error {
			return c.Status(200).JSON(fiber.Map{"result": true})
		})

		req := httptest.NewRequest("GET", "/json", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		entries := store.List()
		require.Len(t, entries, 1)
		assert.NotNil(t, entries[0].Response.Body)
	})

	t.Run("skips admin routes", func(t *testing.T) {
		store := NewStore(100)
		app := fiber.New()
		app.Use(NewRecorderMiddleware(store))
		app.Get("/_inzibat/recorder/entries", func(c *fiber.Ctx) error {
			return c.Status(200).SendString("admin")
		})
		app.Get("/normal", func(c *fiber.Ctx) error {
			return c.Status(200).SendString("ok")
		})

		adminReq := httptest.NewRequest("GET", "/_inzibat/recorder/entries", nil)
		resp, err := app.Test(adminReq, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		normalReq := httptest.NewRequest("GET", "/normal", nil)
		resp2, err := app.Test(normalReq, -1)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, 1, store.Len())
		entries := store.List()
		assert.Equal(t, "/normal", entries[0].Request.Path)
	})

	t.Run("records multiple requests", func(t *testing.T) {
		store := NewStore(100)
		app := fiber.New()
		app.Use(NewRecorderMiddleware(store))
		app.Get("/a", func(c *fiber.Ctx) error {
			return c.Status(200).SendString("a")
		})
		app.Post("/b", func(c *fiber.Ctx) error {
			return c.Status(201).SendString("b")
		})

		req1 := httptest.NewRequest("GET", "/a", nil)
		resp1, err := app.Test(req1, -1)
		require.NoError(t, err)
		defer resp1.Body.Close()

		req2 := httptest.NewRequest("POST", "/b", nil)
		resp2, err := app.Test(req2, -1)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, 2, store.Len())
	})

	t.Run("handles non-JSON body gracefully", func(t *testing.T) {
		store := NewStore(100)
		app := fiber.New()
		app.Use(NewRecorderMiddleware(store))
		app.Post("/text", func(c *fiber.Ctx) error {
			return c.Status(200).SendString("plain response")
		})

		req := httptest.NewRequest("POST", "/text", strings.NewReader("plain text body"))
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		entries := store.List()
		require.Len(t, entries, 1)
		assert.NotNil(t, entries[0].Request.Body)
	})
}

func TestCaptureBody(t *testing.T) {
	t.Run("nil for empty body", func(t *testing.T) {
		result := captureBody(nil)
		assert.Nil(t, result)
	})

	t.Run("nil for zero-length body", func(t *testing.T) {
		result := captureBody([]byte{})
		assert.Nil(t, result)
	})

	t.Run("captures valid JSON", func(t *testing.T) {
		body := []byte(`{"key":"value"}`)
		result := captureBody(body)
		assert.Equal(t, `{"key":"value"}`, string(result))
	})

	t.Run("wraps non-JSON as string", func(t *testing.T) {
		body := []byte("hello world")
		result := captureBody(body)
		assert.Equal(t, `"hello world"`, string(result))
	})

	t.Run("truncates oversized body", func(t *testing.T) {
		body := make([]byte, MaxBodyCaptureBytes+100)
		for i := range body {
			body[i] = 'a'
		}
		result := captureBody(body)
		assert.NotNil(t, result)
	})
}

func BenchmarkMiddleware(b *testing.B) {
	store := NewStore(10000)
	app := fiber.New()
	app.Use(NewRecorderMiddleware(store))
	app.Get("/bench", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	b.ResetTimer()
	for range b.N {
		resp, _ := app.Test(req, -1)
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}
