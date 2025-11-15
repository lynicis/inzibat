package form_builder

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"inzibat/config"
)

func TestHeadersCollector_GetEmptyValue(t *testing.T) {
	t.Run("happy path - returns empty http.Header", func(t *testing.T) {
		collector := &HeadersCollector{}

		result := collector.GetEmptyValue()

		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.NotNil(t, headers)
		assert.Equal(t, 0, len(headers))
	})
}

func TestHeadersCollector_GetSourceTitle(t *testing.T) {
	t.Run("happy path - returns Header Source", func(t *testing.T) {
		collector := &HeadersCollector{}

		result := collector.GetSourceTitle()

		assert.Equal(t, "Header Source", result)
	})
}

func TestHeadersCollector_GetFileFormConfig(t *testing.T) {
	t.Run("happy path - returns correct FilePathFormConfig", func(t *testing.T) {
		collector := &HeadersCollector{}

		result := collector.GetFileFormConfig()

		assert.Equal(t, "filepath", result.Key)
		assert.Equal(t, "Header JSON File Path", result.Title)
		assert.Equal(t, "/path/to/headers.json", result.Placeholder)
	})
}

func TestHeadersCollector_CollectFromFile(t *testing.T) {
	t.Run("happy path - loads headers from valid JSON file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "headers.json")
		jsonContent := `{"X-Custom-Header": "custom-value", "Accept": "application/json"}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		collector := &HeadersCollector{}
		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, "custom-value", headers.Get("X-Custom-Header"))
		assert.Equal(t, "application/json", headers.Get("Accept"))
	})

	t.Run("error path - file does not exist", func(t *testing.T) {
		collector := &HeadersCollector{}
		result, err := collector.CollectFromFile("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("error path - invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `not valid json`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		collector := &HeadersCollector{}
		result, err := collector.CollectFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBodyCollector_GetEmptyValue(t *testing.T) {
	t.Run("happy path - returns config.HttpBody with nil", func(t *testing.T) {
		collector := &BodyCollector{}

		result := collector.GetEmptyValue()

		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Equal(t, 0, len(body))
	})
}

func TestBodyCollector_GetSourceTitle(t *testing.T) {
	t.Run("happy path - returns Body Source", func(t *testing.T) {
		collector := &BodyCollector{}

		result := collector.GetSourceTitle()

		assert.Equal(t, "Body Source", result)
	})
}

func TestBodyCollector_GetFileFormConfig(t *testing.T) {
	t.Run("happy path - returns correct FilePathFormConfig", func(t *testing.T) {
		collector := &BodyCollector{}

		result := collector.GetFileFormConfig()

		assert.Equal(t, "filepath", result.Key)
		assert.Equal(t, "Body JSON File Path", result.Title)
		assert.Equal(t, "/path/to/body.json", result.Placeholder)
	})
}

func TestBodyCollector_CollectFromFile(t *testing.T) {
	t.Run("happy path - loads body from valid JSON file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.json")
		jsonContent := `{"id": 1, "name": "test", "active": true}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		collector := &BodyCollector{}
		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Equal(t, float64(1), body["id"])
		assert.Equal(t, "test", body["name"])
		assert.Equal(t, true, body["active"])
	})

	t.Run("error path - file does not exist", func(t *testing.T) {
		collector := &BodyCollector{}
		result, err := collector.CollectFromFile("/nonexistent/file.json")

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("error path - invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `{invalid json}`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		collector := &BodyCollector{}
		result, err := collector.CollectFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBodyStringCollector_GetEmptyValue(t *testing.T) {
	t.Run("happy path - returns empty string", func(t *testing.T) {
		collector := &BodyStringCollector{}

		result := collector.GetEmptyValue()

		assert.NotNil(t, result)
		str, ok := result.(string)
		assert.True(t, ok)
		assert.Equal(t, "", str)
	})
}

func TestBodyStringCollector_GetSourceTitle(t *testing.T) {
	t.Run("happy path - returns BodyString Source", func(t *testing.T) {
		collector := &BodyStringCollector{}

		result := collector.GetSourceTitle()

		assert.Equal(t, "BodyString Source", result)
	})
}

func TestBodyStringCollector_GetFileFormConfig(t *testing.T) {
	t.Run("happy path - returns correct FilePathFormConfig with FilePathKey", func(t *testing.T) {
		collector := &BodyStringCollector{}

		result := collector.GetFileFormConfig()

		assert.Equal(t, FilePathKey, result.Key)
		assert.Equal(t, "BodyString File Path", result.Title)
		assert.Equal(t, "/path/to/body.txt", result.Placeholder)
	})
}

func TestBodyStringCollector_CollectFromFile(t *testing.T) {
	t.Run("happy path - loads body string from text file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.txt")
		content := "Simple text content"
		err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)

		collector := &BodyStringCollector{}
		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		str, ok := result.(string)
		assert.True(t, ok)
		assert.Equal(t, content, str)
	})

	t.Run("error path - file does not exist", func(t *testing.T) {
		collector := &BodyStringCollector{}
		result, err := collector.CollectFromFile("/nonexistent/file.txt")

		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("happy path - loads empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		collector := &BodyStringCollector{}
		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		str, ok := result.(string)
		assert.True(t, ok)
		assert.Equal(t, "", str)
	})
}

func TestCollectHeaders(t *testing.T) {
	t.Run("happy path - creates HeadersCollector and calls CollectData", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})

	t.Run("error path - CollectData returns error", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})
}

func TestCollectBody(t *testing.T) {
	t.Run("happy path - creates BodyCollector and calls CollectData", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})

	t.Run("happy path - handles nil result", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})

	t.Run("error path - CollectData returns error", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})
}

func TestCollectBodyString(t *testing.T) {
	t.Run("happy path - creates BodyStringCollector and calls CollectData", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})

	t.Run("error path - CollectData returns error", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})
}
