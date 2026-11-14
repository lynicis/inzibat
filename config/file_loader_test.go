package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

	return filePath
}

func TestLoadHeadersFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		filePath := writeTempFile(t, "headers.json",
			`{"Content-Type": "application/json", "Authorization": "Bearer token123"}`)

		headers, err := LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
	})

	t.Run("file not found", func(t *testing.T) {
		headers, err := LoadHeadersFromFile("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		filePath := writeTempFile(t, "invalid.json", `{"key": "value" invalid}`)

		headers, err := LoadHeadersFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("non-map JSON", func(t *testing.T) {
		filePath := writeTempFile(t, "array.json", `[1, 2, 3]`)

		headers, err := LoadHeadersFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("empty object", func(t *testing.T) {
		filePath := writeTempFile(t, "empty.json", `{}`)

		headers, err := LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.Empty(t, headers)
	})

	t.Run("read error", func(t *testing.T) {
		headers, err := LoadHeadersFromFile(t.TempDir())

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to read")
	})
}

func TestLoadBodyFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		filePath := writeTempFile(t, "body.json", `{"message": "Hello", "status": "ok", "count": 42}`)

		body, err := LoadBodyFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "Hello", body["message"])
		assert.Equal(t, "ok", body["status"])
		assert.Equal(t, float64(42), body["count"])
	})

	t.Run("nested JSON object", func(t *testing.T) {
		filePath := writeTempFile(t, "nested.json", `{"outer": {"inner": "value"}}`)

		body, err := LoadBodyFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"inner": "value"}, body["outer"])
	})

	t.Run("file not found", func(t *testing.T) {
		body, err := LoadBodyFromFile("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		filePath := writeTempFile(t, "invalid.json", `not a valid json`)

		body, err := LoadBodyFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("empty file", func(t *testing.T) {
		filePath := writeTempFile(t, "empty.json", "")

		body, err := LoadBodyFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("read error", func(t *testing.T) {
		body, err := LoadBodyFromFile(t.TempDir())

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to read")
	})
}

func TestLoadBodyStringFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		filePath := writeTempFile(t, "body.txt", "Simple text content")

		bodyString, err := LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "Simple text content", bodyString)
	})

	t.Run("JSON content is kept verbatim", func(t *testing.T) {
		filePath := writeTempFile(t, "body.json", `{"message": "success", "status": "ok"}`)

		bodyString, err := LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, `{"message": "success", "status": "ok"}`, bodyString)
	})

	t.Run("file not found", func(t *testing.T) {
		bodyString, err := LoadBodyStringFromFile("/nonexistent/file.txt")

		assert.Error(t, err)
		assert.Equal(t, "", bodyString)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("empty file", func(t *testing.T) {
		filePath := writeTempFile(t, "empty.txt", "")

		bodyString, err := LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("read error", func(t *testing.T) {
		bodyString, err := LoadBodyStringFromFile(t.TempDir())

		assert.Error(t, err)
		assert.Equal(t, "", bodyString)
		assert.Contains(t, err.Error(), "failed to read")
	})
}
