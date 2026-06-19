package recorder

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminApp() (*fiber.App, *Store) {
	store := NewStore(100)
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	RegisterAdminRoutes(app, store)

	return app, store
}

func TestListEntriesHandler(t *testing.T) {
	t.Run("returns empty array when no entries", func(t *testing.T) {
		app, _ := setupAdminApp()
		req := httptest.NewRequest("GET", "/_inzibat/recorder/entries", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "[]", string(body))
	})

	t.Run("returns recorded entries", func(t *testing.T) {
		app, store := setupAdminApp()
		store.Add(RecordedEntry{
			ID:        "test-1",
			Timestamp: time.Now(),
			Request: RecordedRequest{
				Method: "GET",
				Path:   "/hello",
			},
			Response: RecordedResponse{
				StatusCode: 200,
			},
		})

		req := httptest.NewRequest("GET", "/_inzibat/recorder/entries", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var entries []RecordedEntry
		err = json.Unmarshal(body, &entries)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, "test-1", entries[0].ID)
		assert.Equal(t, "GET", entries[0].Request.Method)
	})
}

func TestGetSessionHandler(t *testing.T) {
	t.Run("returns session with metadata", func(t *testing.T) {
		app, store := setupAdminApp()
		store.Add(RecordedEntry{
			ID:        "s-1",
			Timestamp: time.Now(),
			Request:   RecordedRequest{Method: "POST", Path: "/data"},
			Response:  RecordedResponse{StatusCode: 201},
		})

		req := httptest.NewRequest("GET", "/_inzibat/recorder/session", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var session RecordedSession
		err = json.Unmarshal(body, &session)
		require.NoError(t, err)
		assert.Equal(t, 1, session.EntryCount)
		assert.False(t, session.StartedAt.IsZero())
		require.Len(t, session.Entries, 1)
	})
}

func TestClearEntriesHandler(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		app, store := setupAdminApp()
		store.Add(RecordedEntry{
			ID:       "c-1",
			Request:  RecordedRequest{Method: "GET", Path: "/test"},
			Response: RecordedResponse{StatusCode: 200},
		})
		store.Add(RecordedEntry{
			ID:       "c-2",
			Request:  RecordedRequest{Method: "POST", Path: "/test"},
			Response: RecordedResponse{StatusCode: 201},
		})
		assert.Equal(t, 2, store.Len())

		req := httptest.NewRequest("POST", "/_inzibat/recorder/clear", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 0, store.Len())

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "cleared")
	})
}
