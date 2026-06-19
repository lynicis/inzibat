package recorder

import (
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToInzibatConfig(t *testing.T) {
	t.Run("converts basic session", func(t *testing.T) {
		session := RecordedSession{
			StartedAt: time.Now(),
			Entries: []RecordedEntry{
				{
					Request: RecordedRequest{
						Method: "GET",
						Path:   "/users",
					},
					Response: RecordedResponse{
						StatusCode: 200,
						Body:       json.RawMessage(`{"users":[]}`),
					},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 9090)
		assert.Equal(t, 9090, cfg.ServerPort)
		require.Len(t, cfg.Routes, 1)
		assert.Equal(t, "GET", cfg.Routes[0].Method)
		assert.Equal(t, "/users", cfg.Routes[0].Path)
		assert.NotNil(t, cfg.Routes[0].FakeResponse)
		assert.Equal(t, 200, cfg.Routes[0].FakeResponse.StatusCode)
	})

	t.Run("deduplicates by method and path keeping last", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request: RecordedRequest{Method: "GET", Path: "/api"},
					Response: RecordedResponse{
						StatusCode: 200,
						Body:       json.RawMessage(`{"version":"v1"}`),
					},
				},
				{
					Request: RecordedRequest{Method: "GET", Path: "/api"},
					Response: RecordedResponse{
						StatusCode: 200,
						Body:       json.RawMessage(`{"version":"v2"}`),
					},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 8080)
		require.Len(t, cfg.Routes, 1)
		assert.Equal(t, "v2", cfg.Routes[0].FakeResponse.Body["version"])
	})

	t.Run("preserves different methods on same path", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request:  RecordedRequest{Method: "GET", Path: "/items"},
					Response: RecordedResponse{StatusCode: 200},
				},
				{
					Request:  RecordedRequest{Method: "POST", Path: "/items"},
					Response: RecordedResponse{StatusCode: 201},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 8080)
		require.Len(t, cfg.Routes, 2)
		assert.Equal(t, "GET", cfg.Routes[0].Method)
		assert.Equal(t, "POST", cfg.Routes[1].Method)
	})

	t.Run("handles string body as bodyString", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request:  RecordedRequest{Method: "GET", Path: "/text"},
					Response: RecordedResponse{StatusCode: 200, Body: json.RawMessage(`"hello world"`)},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 8080)
		require.Len(t, cfg.Routes, 1)
		assert.Equal(t, "hello world", cfg.Routes[0].FakeResponse.BodyString)
		assert.Nil(t, cfg.Routes[0].FakeResponse.Body)
	})

	t.Run("handles empty body", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request:  RecordedRequest{Method: "DELETE", Path: "/item/1"},
					Response: RecordedResponse{StatusCode: 204},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 8080)
		require.Len(t, cfg.Routes, 1)
		assert.Nil(t, cfg.Routes[0].FakeResponse.Body)
		assert.Empty(t, cfg.Routes[0].FakeResponse.BodyString)
	})

	t.Run("handles response headers", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request: RecordedRequest{Method: "GET", Path: "/with-headers"},
					Response: RecordedResponse{
						StatusCode: 200,
						Headers: map[string][]string{
							"Content-Type": {"application/json"},
							"X-Custom":     {"custom-value"},
						},
						Body: json.RawMessage(`{"ok":true}`),
					},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 8080)
		require.Len(t, cfg.Routes, 1)
		assert.Equal(t, []string{"application/json"}, cfg.Routes[0].FakeResponse.Headers["Content-Type"])
		assert.Equal(t, []string{"custom-value"}, cfg.Routes[0].FakeResponse.Headers["X-Custom"])
	})

	t.Run("uses default port when zero", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request:  RecordedRequest{Method: "GET", Path: "/"},
					Response: RecordedResponse{StatusCode: 200},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 0)
		assert.Equal(t, 8080, cfg.ServerPort)
	})

	t.Run("handles empty session", func(t *testing.T) {
		session := RecordedSession{Entries: []RecordedEntry{}}

		cfg := ConvertToInzibatConfig(session, 8080)
		assert.Equal(t, 8080, cfg.ServerPort)
		assert.Empty(t, cfg.Routes)
		assert.NotNil(t, cfg.Routes)
	})

	t.Run("preserves route order", func(t *testing.T) {
		session := RecordedSession{
			Entries: []RecordedEntry{
				{
					Request:  RecordedRequest{Method: "GET", Path: "/c"},
					Response: RecordedResponse{StatusCode: 200},
				},
				{
					Request:  RecordedRequest{Method: "GET", Path: "/a"},
					Response: RecordedResponse{StatusCode: 200},
				},
				{
					Request:  RecordedRequest{Method: "GET", Path: "/b"},
					Response: RecordedResponse{StatusCode: 200},
				},
			},
		}

		cfg := ConvertToInzibatConfig(session, 8080)
		require.Len(t, cfg.Routes, 3)
		assert.Equal(t, "/c", cfg.Routes[0].Path)
		assert.Equal(t, "/a", cfg.Routes[1].Path)
		assert.Equal(t, "/b", cfg.Routes[2].Path)
	})
}

func TestBuildFakeResponse(t *testing.T) {
	t.Run("JSON object body goes to Body field", func(t *testing.T) {
		resp := RecordedResponse{
			StatusCode: 200,
			Body:       json.RawMessage(`{"name":"test","count":42}`),
		}

		fake := buildFakeResponse(resp)
		assert.Equal(t, 200, fake.StatusCode)
		assert.Equal(t, "test", fake.Body["name"])
		assert.Equal(t, float64(42), fake.Body["count"])
	})

	t.Run("JSON array falls back to bodyString", func(t *testing.T) {
		resp := RecordedResponse{
			StatusCode: 200,
			Body:       json.RawMessage(`[1,2,3]`),
		}

		fake := buildFakeResponse(resp)
		assert.Nil(t, fake.Body)
		assert.Equal(t, "[1,2,3]", fake.BodyString)
	})
}

func TestConvertHeaders(t *testing.T) {
	t.Run("nil for empty headers", func(t *testing.T) {
		result := convertHeaders(nil)
		assert.Nil(t, result)
	})

	t.Run("converts headers", func(t *testing.T) {
		headers := map[string][]string{
			"Content-Type": {"application/json"},
		}

		result := convertHeaders(headers)
		assert.Equal(t, []string{"application/json"}, result["Content-Type"])
	})
}
