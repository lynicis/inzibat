package form_builder

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func TestHeadersCollector(t *testing.T) {
	t.Run("happy path - GetEmptyValue returns empty header", func(t *testing.T) {
		collector := &HeadersCollector{}

		emptyValue := collector.GetEmptyValue()

		assert.NotNil(t, emptyValue)
		headers, ok := emptyValue.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("happy path - GetSourceTitle returns correct title", func(t *testing.T) {
		collector := &HeadersCollector{}

		title := collector.GetSourceTitle()

		assert.Equal(t, "Header Source", title)
	})

	t.Run("happy path - GetFileFormConfig returns correct config", func(t *testing.T) {
		collector := &HeadersCollector{}

		config := collector.GetFileFormConfig()

		assert.Equal(t, "filepath", config.Key)
		assert.Equal(t, "Header JSON File Path", config.Title)
		assert.Equal(t, "/path/to/headers.json", config.Placeholder)
	})

	t.Run("happy path - CollectFromFile loads headers", func(t *testing.T) {
		collector := &HeadersCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "headers.json")
		headersData := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token",
		}
		data, err := json.Marshal(headersData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token", headers.Get("Authorization"))
	})
}

func TestBodyCollector(t *testing.T) {
	t.Run("happy path - GetEmptyValue returns nil HttpBody", func(t *testing.T) {
		collector := &BodyCollector{}

		emptyValue := collector.GetEmptyValue()

		body, ok := emptyValue.(config.HttpBody)
		assert.True(t, ok)
		if body != nil {
			assert.Equal(t, 0, len(body))
		}
	})

	t.Run("happy path - GetSourceTitle returns correct title", func(t *testing.T) {
		collector := &BodyCollector{}

		title := collector.GetSourceTitle()

		assert.Equal(t, "Body Source", title)
	})

	t.Run("happy path - GetFileFormConfig returns correct config", func(t *testing.T) {
		collector := &BodyCollector{}

		config := collector.GetFileFormConfig()

		assert.Equal(t, "filepath", config.Key)
		assert.Equal(t, "Body JSON File Path", config.Title)
		assert.Equal(t, "/path/to/body.json", config.Placeholder)
	})

	t.Run("happy path - CollectFromFile loads body", func(t *testing.T) {
		collector := &BodyCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.json")
		bodyData := config.HttpBody{
			"message": "success",
			"code":    float64(200),
		}
		data, err := json.Marshal(bodyData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, float64(200), body["code"])
	})
}

func TestBodyStringCollector(t *testing.T) {
	t.Run("happy path - GetEmptyValue returns empty string", func(t *testing.T) {
		collector := &BodyStringCollector{}

		emptyValue := collector.GetEmptyValue()

		assert.Equal(t, "", emptyValue)
	})

	t.Run("happy path - GetSourceTitle returns correct title", func(t *testing.T) {
		collector := &BodyStringCollector{}

		title := collector.GetSourceTitle()

		assert.Equal(t, "BodyString Source", title)
	})

	t.Run("happy path - GetFileFormConfig returns correct config", func(t *testing.T) {
		collector := &BodyStringCollector{}

		config := collector.GetFileFormConfig()

		assert.Equal(t, FilePathKey, config.Key)
		assert.Equal(t, "BodyString File Path", config.Title)
		assert.Equal(t, "/path/to/body.txt", config.Placeholder)
	})

	t.Run("happy path - CollectFromFile loads body string", func(t *testing.T) {
		collector := &BodyStringCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.txt")
		expectedContent := `{"message": "success"}`
		err := os.WriteFile(filePath, []byte(expectedContent), 0644)
		require.NoError(t, err)

		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, result)
	})
}

func TestCollectHeaders(t *testing.T) {
	t.Run("happy path - returns headers collector interface", func(t *testing.T) {
		collector := &HeadersCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})
}

func TestCollectBody(t *testing.T) {
	t.Run("happy path - returns body collector interface", func(t *testing.T) {
		collector := &BodyCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})
}

func TestCollectBodyString(t *testing.T) {
	t.Run("happy path - returns body string collector interface", func(t *testing.T) {
		collector := &BodyStringCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})
}
