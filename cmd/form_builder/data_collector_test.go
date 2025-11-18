package form_builder

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lynicis/inzibat/config"
)

func TestCollectHeaders(t *testing.T) {
	t.Run("happy path - CollectHeaders function exists and returns correct type", func(t *testing.T) {
		var _ func() (http.Header, error) = CollectHeaders

		assert.NotNil(t, CollectHeaders)
	})
}

func TestCollectBody(t *testing.T) {
	t.Run("happy path - CollectBody function exists and returns correct type", func(t *testing.T) {
		var _ func() (config.HttpBody, error) = CollectBody

		assert.NotNil(t, CollectBody)
	})
}

func TestCollectBodyString(t *testing.T) {
	t.Run("happy path - CollectBodyString function exists and returns correct type", func(t *testing.T) {
		var _ func() (string, error) = CollectBodyString

		assert.NotNil(t, CollectBodyString)
	})
}

func TestCollectHeadersInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - sourceFormRunner.Run() returns error", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		expectedError := errors.New("form run error")

		mockSourceForm.EXPECT().Run().Return(expectedError)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, headers)
	})

	t.Run("happy path - SourceSkip returns empty headers", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceSkip)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("error path - SourceFile but filePathFormRunner.Run() returns error", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		expectedError := errors.New("file path form error")

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(expectedError)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, headers)
	})

	t.Run("happy path - SourceFile loads headers from file successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "headers.json")
		jsonContent := `{"Content-Type": "application/json", "Authorization": "Bearer token"}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token", headers.Get("Authorization"))
	})

	t.Run("error path - SourceFile but file does not exist", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		nonExistentPath := "/nonexistent/headers.json"

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(nonExistentPath)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, headers)
	})

	t.Run("happy path - SourceForm calls CollectHeadersFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectHeadersFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceForm)

		_, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})

	t.Run("error path - SourceFile but invalid JSON in headers file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid_headers.json")
		invalidJSON := `{invalid json}`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("happy path - SourceFile loads empty headers file", func(t *testing.T) {

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty_headers.json")
		emptyJSON := `{}`
		err := os.WriteFile(filePath, []byte(emptyJSON), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("happy path - SourceFile loads headers with multiple values", func(t *testing.T) {

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "multi_headers.json")
		jsonContent := `{"Accept": "application/json", "X-Custom": "value1", "Authorization": "Bearer token123"}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		headers, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, "application/json", headers.Get("Accept"))
		assert.Equal(t, "value1", headers.Get("X-Custom"))
		assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
	})

	t.Run("happy path - unknown source falls through to CollectHeadersFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectHeadersFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return("unknown_source")

		_, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})

	t.Run("happy path - empty source falls through to CollectHeadersFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectHeadersFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return("")

		_, err := collectHeadersInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})
}

func TestCollectBodyInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - sourceFormRunner.Run() returns error", func(t *testing.T) {

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		expectedError := errors.New("form run error")

		mockSourceForm.EXPECT().Run().Return(expectedError)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, body)
	})

	t.Run("happy path - SourceSkip returns nil body", func(t *testing.T) {

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceSkip)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.Nil(t, body)
	})

	t.Run("error path - SourceFile but filePathFormRunner.Run() returns error", func(t *testing.T) {

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		expectedError := errors.New("file path form error")

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(expectedError)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, body)
	})

	t.Run("happy path - SourceFile loads body from file successfully", func(t *testing.T) {

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.json")
		jsonContent := `{"message": "success", "id": 1, "active": true}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, body)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, float64(1), body["id"])
		assert.Equal(t, true, body["active"])
	})

	t.Run("error path - SourceFile but file does not exist", func(t *testing.T) {

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		nonExistentPath := "/nonexistent/body.json"

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(nonExistentPath)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, body)
	})

	t.Run("error path - SourceFile but invalid JSON", func(t *testing.T) {

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.json")
		invalidJSON := `{invalid json}`
		err := os.WriteFile(filePath, []byte(invalidJSON), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, body)
	})

	t.Run("happy path - SourceForm calls CollectBodyFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectBodyFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceForm)

		_, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})

	t.Run("happy path - SourceFile loads empty body file", func(t *testing.T) {

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty_body.json")
		emptyJSON := `{}`
		err := os.WriteFile(filePath, []byte(emptyJSON), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, body)
		assert.Equal(t, 0, len(body))
	})

	t.Run("happy path - SourceFile loads body with nested objects", func(t *testing.T) {

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "nested_body.json")
		jsonContent := `{"user": {"id": 1, "name": "test"}, "tags": ["tag1", "tag2"]}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString("filepath").Return(filePath)

		body, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.NotNil(t, body)
		assert.NotNil(t, body["user"])
		assert.NotNil(t, body["tags"])
	})

	t.Run("happy path - unknown source falls through to CollectBodyFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectBodyFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return("unknown_source")

		_, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})

	t.Run("happy path - empty source falls through to CollectBodyFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectBodyFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return("")

		_, err := collectBodyInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})
}

func TestCollectBodyStringInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - sourceFormRunner.Run() returns error", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		expectedError := errors.New("form run error")

		mockSourceForm.EXPECT().Run().Return(expectedError)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("happy path - SourceSkip returns empty string", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceSkip)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("error path - SourceFile but filePathFormRunner.Run() returns error", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		expectedError := errors.New("file path form error")

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(expectedError)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("happy path - SourceFile loads body string from file successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "body.txt")
		content := "Simple text content for body string"
		err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString(FilePathKey).Return(filePath)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.Equal(t, content, bodyString)
	})

	t.Run("error path - SourceFile but file does not exist", func(t *testing.T) {
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)
		nonExistentPath := "/nonexistent/body.txt"

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString(FilePathKey).Return(nonExistentPath)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.Error(t, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("happy path - SourceFile loads empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString(FilePathKey).Return(filePath)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.Equal(t, "", bodyString)
	})

	t.Run("happy path - SourceForm calls CollectBodyStringFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectBodyStringFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceForm)

		_, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})

	t.Run("happy path - SourceFile loads multiline body string", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "multiline.txt")
		content := "Line 1\nLine 2\nLine 3\nWith special chars: {}\""
		err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString(FilePathKey).Return(filePath)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.Equal(t, content, bodyString)
	})

	t.Run("happy path - SourceFile loads JSON string as body string", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "json_string.txt")
		jsonContent := `{"message": "success", "id": 1}`
		err := os.WriteFile(filePath, []byte(jsonContent), 0644)
		assert.NoError(t, err)

		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return(SourceFile)
		mockFilePathForm.EXPECT().Run().Return(nil)
		mockFilePathForm.EXPECT().GetString(FilePathKey).Return(filePath)

		bodyString, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		assert.NoError(t, err)
		assert.Equal(t, jsonContent, bodyString)
	})

	t.Run("happy path - unknown source falls through to CollectBodyStringFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectBodyStringFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return("unknown_source")

		_, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})

	t.Run("happy path - empty source falls through to CollectBodyStringFromForm", func(t *testing.T) {
		t.Skip("Skipping interactive form test - CollectBodyStringFromForm requires TTY and will hang in non-interactive environments")
		mockSourceForm := NewMockFormRunner(ctrl)
		mockFilePathForm := NewMockFormRunner(ctrl)

		mockSourceForm.EXPECT().Run().Return(nil)
		mockSourceForm.EXPECT().GetString(SourceKey).Return("")

		_, err := collectBodyStringInternal(mockSourceForm, mockFilePathForm)

		_ = err
	})
}
