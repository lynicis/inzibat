package recorder

import (
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const adminPathPrefix = "/_inzibat/"

// NewRecorderMiddleware creates a Fiber middleware that captures request/response pairs
// into the provided Store. Admin routes (/_inzibat/*) are skipped.
func NewRecorderMiddleware(store *Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if strings.HasPrefix(ctx.Path(), adminPathPrefix) {
			return ctx.Next()
		}

		start := time.Now()

		reqBody := captureBody(ctx.Body())
		reqHeaders := captureRequestHeaders(ctx)

		err := ctx.Next()

		duration := time.Since(start).Milliseconds()
		respBody := captureBody(ctx.Response().Body())
		respHeaders := captureResponseHeaders(ctx)

		entry := RecordedEntry{
			ID:        uuid.NewString(),
			Timestamp: start,
			Request: RecordedRequest{
				Method:  ctx.Method(),
				Path:    ctx.Path(),
				Headers: reqHeaders,
				Body:    reqBody,
			},
			Response: RecordedResponse{
				StatusCode: ctx.Response().StatusCode(),
				Headers:    respHeaders,
				Body:       respBody,
			},
			DurationMs: duration,
		}

		store.Add(entry)

		return err
	}
}

func captureBody(body []byte) json.RawMessage {
	if len(body) == 0 {
		return nil
	}

	if len(body) > MaxBodyCaptureBytes {
		body = body[:MaxBodyCaptureBytes]
	}

	// If it's valid JSON, store as-is
	if json.Valid(body) {
		result := make(json.RawMessage, len(body))
		copy(result, body)
		return result
	}

	// Non-JSON body: wrap as a JSON string
	encoded, err := json.Marshal(string(body))
	if err != nil {
		return nil
	}

	return encoded
}

func captureRequestHeaders(c *fiber.Ctx) map[string][]string {
	headers := map[string][]string{}

	for key, value := range c.Request().Header.All() {
		k := string(key)
		headers[k] = append(headers[k], string(value))
	}

	if len(headers) == 0 {
		return nil
	}

	return headers
}

func captureResponseHeaders(c *fiber.Ctx) map[string][]string {
	headers := map[string][]string{}

	for key, value := range c.Response().Header.All() {
		k := string(key)
		headers[k] = append(headers[k], string(value))
	}

	if len(headers) == 0 {
		return nil
	}

	return headers
}
