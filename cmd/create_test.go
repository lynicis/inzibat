package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func TestLoadHeadersFromFile(t *testing.T) {
	t.Run("happy path - valid JSON file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "headers.json")
		headersData := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
		}
		data, err := json.Marshal(headersData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		headers, err := config.LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
	})

	t.Run("error path - file does not exist", func(t *testing.T) {
		nonExistentPath := "/non/existent/file.json"

		headers, err := config.LoadHeadersFromFile(nonExistentPath)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("error path - invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(filePath, []byte("invalid json content"), 0644)
		require.NoError(t, err)

		headers, err := config.LoadHeadersFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("error path - empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty.json")
		err := os.WriteFile(filePath, []byte("{}"), 0644)
		require.NoError(t, err)

		headers, err := config.LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, 0, len(headers))
	})
}

func TestLoadBodyFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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

		body, err := config.LoadBodyFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, body)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, float64(200), body["code"])
	})

	t.Run("file does not exist", func(t *testing.T) {
		nonExistentPath := "/non/existent/body.json"

		body, err := config.LoadBodyFromFile(nonExistentPath)

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(filePath, []byte("not a valid json"), 0644)
		require.NoError(t, err)

		body, err := config.LoadBodyFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})
}

func TestLoadBodyStringFromFile(t *testing.T) {
	t.Run("happy path - valid text file", func(t *testing.T) {

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.txt")
		expectedContent := `{"message": "success", "status": "ok"}`
		err := os.WriteFile(filePath, []byte(expectedContent), 0644)
		require.NoError(t, err)

		bodyString, err := config.LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, bodyString)
	})

	t.Run("error path - file does not exist", func(t *testing.T) {

		nonExistentPath := "/non/existent/body.txt"

		bodyString, err := config.LoadBodyStringFromFile(nonExistentPath)

		// Assert

		assert.Empty(t, bodyString)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("happy path - empty file", func(t *testing.T) {

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty.txt")
		err := os.WriteFile(filePath, []byte(""), 0644)
		require.NoError(t, err)

		bodyString, err := config.LoadBodyStringFromFile(filePath)

		// Assert

		assert.Empty(t, bodyString)
	})
}

func TestCreateRouteForm(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// Arrange & Act
		form := createRouteForm()

		// Assert
		assert.NotNil(t, form)
	})
}
