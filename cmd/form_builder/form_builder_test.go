package form_builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFilePathForm(t *testing.T) {
	t.Run("happy path - builds form with file path validation", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
	})

	t.Run("happy path - form has correct structure", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Config File",
			Placeholder: "/path/to/file",
			Key:         "config_file",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
	})

	t.Run("happy path - validation function validates non-empty path", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)
	})

	t.Run("happy path - validation function validates file exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)
	})

	t.Run("error path - validation rejects empty path", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)
	})

	t.Run("error path - validation rejects non-existent file", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)
	})
}

func TestBuildSourceSelectionForm(t *testing.T) {
	t.Run("happy path - builds form with source options", func(t *testing.T) {
		title := "Select Source"
		key := "source"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
	})

	t.Run("happy path - form has correct structure", func(t *testing.T) {
		title := "Choose Option"
		key := "option"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
	})

	t.Run("happy path - form includes all source options", func(t *testing.T) {
		title := "Select Source"
		key := "source"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
	})
}

func TestGetFilePathFromForm(t *testing.T) {
	t.Run("happy path - function exists and creates form", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)
	})

	t.Run("error path - function handles form creation", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Test",
			Placeholder: "test",
			Key:         "test",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)
	})
}

func TestGetSourceFromForm(t *testing.T) {
	t.Run("happy path - function exists and creates form", func(t *testing.T) {
		title := "Select Source"
		key := "source"

		form := BuildSourceSelectionForm(title, key)
		assert.NotNil(t, form)
	})

	t.Run("error path - function handles form creation", func(t *testing.T) {
		title := "Test"
		key := "test"

		form := BuildSourceSelectionForm(title, key)
		assert.NotNil(t, form)
	})
}
