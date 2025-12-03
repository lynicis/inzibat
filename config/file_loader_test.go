package config

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadersLoader_Load(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "headers.json")
		jsonContent := `{"Content-Type": "application/json", "Authorization": "Bearer token123"}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		loader := &HeadersLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
	})

	t.Run("file not found", func(t *testing.T) {
		loader := &HeadersLoader{}

		result, err := loader.Load("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `{"key": "value" invalid}`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		loader := &HeadersLoader{}

		result, err := loader.Load(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.json")
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		loader := &HeadersLoader{}

		result, err := loader.Load(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("non-map JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "array.json")
		jsonContent := `["value1", "value2"]`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		loader := &HeadersLoader{}

		result, err := loader.Load(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})
}

func TestBodyLoader_Load(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.json")
		jsonContent := `{"message": "Hello", "status": "ok", "count": 42}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		loader := &BodyLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(HttpBody)
		assert.True(t, ok)
		assert.Equal(t, "Hello", body["message"])
		assert.Equal(t, "ok", body["status"])
		assert.Equal(t, float64(42), body["count"])
	})

	t.Run("file not found", func(t *testing.T) {
		loader := &BodyLoader{}

		result, err := loader.Load("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `{"key": "value" invalid}`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		loader := &BodyLoader{}

		result, err := loader.Load(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.json")
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		loader := &BodyLoader{}

		result, err := loader.Load(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("nested JSON object", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "nested.json")
		jsonContent := `{"user": {"name": "John", "age": 30}, "tags": ["admin", "user"]}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		loader := &BodyLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(HttpBody)
		assert.True(t, ok)
		assert.NotNil(t, body["user"])
		assert.NotNil(t, body["tags"])
	})
}

func TestBodyStringLoader_Load(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.txt")
		content := "This is a plain text body"
		err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)

		loader := &BodyStringLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.Equal(t, content, result)
		assert.IsType(t, "", result)
	})

	t.Run("file not found", func(t *testing.T) {
		loader := &BodyStringLoader{}

		result, err := loader.Load("/nonexistent/file.txt")

		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		loader := &BodyStringLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("multiline content", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "multiline.txt")
		content := "Line 1\nLine 2\nLine 3"
		err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)

		loader := &BodyStringLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.Equal(t, content, result)
	})

	t.Run("JSON content as string", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "json.txt")
		jsonContent := `{"key": "value"}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		loader := &BodyStringLoader{}

		result, err := loader.Load(filePath)

		assert.NoError(t, err)
		assert.Equal(t, jsonContent, result)
	})
}

func TestLoadHeadersFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "headers.json")
		jsonContent := `{"X-Custom-Header": "custom-value", "Accept": "application/json"}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		headers, err := LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, "custom-value", headers.Get("X-Custom-Header"))
		assert.Equal(t, "application/json", headers.Get("Accept"))
	})

	t.Run("file not found", func(t *testing.T) {
		headers, err := LoadHeadersFromFile("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, headers)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `not valid json`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		headers, err := LoadHeadersFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, headers)
	})
}

func TestLoadBodyFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.json")
		jsonContent := `{"id": 1, "name": "test", "active": true}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		body, err := LoadBodyFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, body)
		assert.Equal(t, float64(1), body["id"])
		assert.Equal(t, "test", body["name"])
		assert.Equal(t, true, body["active"])
	})

	t.Run("file not found", func(t *testing.T) {
		body, err := LoadBodyFromFile("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, body)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `{invalid json}`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		body, err := LoadBodyFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, body)
	})
}

func TestLoadBodyStringFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.txt")
		content := "Simple text content"
		err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)

		bodyString, err := LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, content, bodyString)
	})

	t.Run("file not found", func(t *testing.T) {
		bodyString, err := LoadBodyStringFromFile("/nonexistent/file.txt")

		assert.Error(t, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		bodyString, err := LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "", bodyString)
	})
}
