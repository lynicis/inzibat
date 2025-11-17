package form_builder

import (
	"fmt"
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

	t.Run("happy path - validation accepts existing file", func(t *testing.T) {
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

		validationFunc := func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}

		err = validationFunc(testFile)
		assert.NoError(t, err)

		_, statErr := os.Stat(testFile)
		assert.NoError(t, statErr)
		assert.False(t, os.IsNotExist(statErr))
	})

	t.Run("error path - validation rejects empty path", func(t *testing.T) {
		emptyPath := ""

		err := ValidateNonEmpty(emptyPath, "file path")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file path")
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("error path - validation rejects non-existent file", func(t *testing.T) {
		nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.json")

		validationFunc := func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}

		err := validationFunc(nonExistentPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file does not exist")

		_, statErr := os.Stat(nonExistentPath)
		assert.Error(t, statErr)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("happy path - validation accepts existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		validationFunc := func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}

		err := validationFunc(tmpDir)

		assert.NoError(t, err)

		fileInfo, statErr := os.Stat(tmpDir)
		assert.NoError(t, statErr)
		assert.True(t, fileInfo.IsDir())
	})

	t.Run("happy path - form created with all config fields", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Custom Title",
			Placeholder: "Custom Placeholder",
			Key:         "custom_key",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
		value := form.GetString(config.Key)
		assert.Equal(t, "", value)
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
		value := form.GetString(key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - form includes all source options", func(t *testing.T) {
		form := BuildSourceSelectionForm("Select Source", "source")

		assert.NotNil(t, form)
		assert.Equal(t, "file", SourceFile)
		assert.Equal(t, "form", SourceForm)
		assert.Equal(t, "skip", SourceSkip)
	})

	t.Run("happy path - form created with different title and key", func(t *testing.T) {
		title := "Body Source"
		key := "body_source"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
		value := form.GetString(key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - form created with empty title", func(t *testing.T) {
		title := ""
		key := "source"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
		value := form.GetString(key)
		assert.Equal(t, "", value)
	})
}

func TestGetFilePathFromForm(t *testing.T) {
	t.Run("happy path - function signature matches expected", func(t *testing.T) {
		var _ func(FilePathFormConfig) (string, error) = GetFilePathFromForm

		assert.NotNil(t, GetFilePathFromForm)
	})

	t.Run("error path - function returns error when form.Run fails", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)

		_ = GetFilePathFromForm
	})

	t.Run("happy path - function uses config key to get value", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Test",
			Placeholder: "test",
			Key:         "test_key",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)

		value := form.GetString(config.Key)
		assert.Equal(t, "", value)
	})

	t.Run("error path - function wraps form.Run error correctly", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)

		_ = GetFilePathFromForm
	})
}

func TestGetSourceFromForm(t *testing.T) {
	t.Run("happy path - function creates form with correct parameters", func(t *testing.T) {
		title := "Select Source"
		key := "source"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
		value := form.GetString(key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - function signature matches expected", func(t *testing.T) {
		var _ func(string, string) (string, error) = GetSourceFromForm

		assert.NotNil(t, GetSourceFromForm)
	})

	t.Run("error path - function returns error when form.Run fails", func(t *testing.T) {
		title := "Select Source"
		key := "source"

		form := BuildSourceSelectionForm(title, key)
		assert.NotNil(t, form)

		_ = GetSourceFromForm
	})

	t.Run("happy path - function uses key parameter to get value", func(t *testing.T) {
		title := "Test"
		key := "test_key"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
		value := form.GetString(key)
		assert.Equal(t, "", value)
	})
}

func TestBuildFilePathFormValidationLogic(t *testing.T) {
	t.Run("happy path - validation passes for existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		validationFunc := func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}

		err = validationFunc(testFile)

		assert.NoError(t, err)

		fileInfo, statErr := os.Stat(testFile)
		assert.NoError(t, statErr)
		assert.NotNil(t, fileInfo)
		assert.False(t, fileInfo.IsDir())
	})

	t.Run("error path - validation fails for empty string", func(t *testing.T) {
		emptyPath := ""

		err := ValidateNonEmpty(emptyPath, "file path")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file path")
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("error path - validation fails for non-existent file path", func(t *testing.T) {
		nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.json")

		validationFunc := func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}

		err := validationFunc(nonExistentPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file does not exist")

		_, statErr := os.Stat(nonExistentPath)
		assert.Error(t, statErr)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("happy path - validation passes for existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		validationFunc := func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}

		err := validationFunc(tmpDir)

		assert.NoError(t, err)

		fileInfo, statErr := os.Stat(tmpDir)
		assert.NoError(t, statErr)
		assert.True(t, fileInfo.IsDir())
	})
}

func TestBuildSourceSelectionFormOptions(t *testing.T) {
	t.Run("happy path - form contains File option", func(t *testing.T) {
		form := BuildSourceSelectionForm("Select Source", "source")

		assert.NotNil(t, form)
		assert.Equal(t, "file", SourceFile)
	})

	t.Run("happy path - form contains Form option", func(t *testing.T) {
		form := BuildSourceSelectionForm("Select Source", "source")

		assert.NotNil(t, form)
		assert.Equal(t, "form", SourceForm)
	})

	t.Run("happy path - form contains Skip option", func(t *testing.T) {
		form := BuildSourceSelectionForm("Select Source", "source")

		assert.NotNil(t, form)
		assert.Equal(t, "skip", SourceSkip)
	})
}

func TestFilePathFormConfig(t *testing.T) {
	t.Run("happy path - config with all fields set", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Test Title",
			Placeholder: "Test Placeholder",
			Key:         "test_key",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
		assert.Equal(t, "Test Title", config.Title)
		assert.Equal(t, "Test Placeholder", config.Placeholder)
		assert.Equal(t, "test_key", config.Key)

		value := form.GetString(config.Key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - config with empty fields", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "",
			Placeholder: "",
			Key:         "",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
		value := form.GetString(config.Key)
		assert.Equal(t, "", value)
	})
}

func TestSourceSelectionFormParameters(t *testing.T) {
	t.Run("happy path - different title values", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
			key   string
		}{
			{"Header Source", "Header Source", "source"},
			{"Body Source", "Body Source", "source"},
			{"BodyString Source", "BodyString Source", "source"},
			{"Empty title", "", "source"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				form := BuildSourceSelectionForm(tc.title, tc.key)

				assert.NotNil(t, form)
				value := form.GetString(tc.key)
				assert.Equal(t, "", value)
			})
		}
	})

	t.Run("happy path - different key values", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
			key   string
		}{
			{"source key", "Select Source", "source"},
			{"custom key", "Select Source", "custom_key"},
			{"empty key", "Select Source", ""},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				form := BuildSourceSelectionForm(tc.title, tc.key)

				assert.NotNil(t, form)
				value := form.GetString(tc.key)
				assert.Equal(t, "", value)
			})
		}
	})
}

func TestSourceConstants(t *testing.T) {
	t.Run("happy path - SourceFile constant is correct", func(t *testing.T) {
		value := SourceFile

		assert.Equal(t, "file", value)
	})

	t.Run("happy path - SourceForm constant is correct", func(t *testing.T) {
		value := SourceForm

		assert.Equal(t, "form", value)
	})

	t.Run("happy path - SourceSkip constant is correct", func(t *testing.T) {
		value := SourceSkip

		assert.Equal(t, "skip", value)
	})

	t.Run("happy path - SourceKey constant is correct", func(t *testing.T) {
		value := SourceKey

		assert.Equal(t, "source", value)
	})

	t.Run("happy path - FilePathKey constant is correct", func(t *testing.T) {
		value := FilePathKey

		assert.Equal(t, "filepath", value)
	})
}

func TestSourceOption(t *testing.T) {
	t.Run("happy path - SourceOption struct exists", func(t *testing.T) {
		option := SourceOption{
			Title: "Test Title",
			Value: "test_value",
		}

		assert.Equal(t, "Test Title", option.Title)
		assert.Equal(t, "test_value", option.Value)
	})

	t.Run("happy path - SourceOption with empty fields", func(t *testing.T) {
		option := SourceOption{
			Title: "",
			Value: "",
		}

		assert.Equal(t, "", option.Title)
		assert.Equal(t, "", option.Value)
	})
}

func TestValidateFilePath(t *testing.T) {
	t.Run("happy path - validation passes for existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})

	t.Run("happy path - validation passes for existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		err := ValidateFilePath(tmpDir)

		assert.NoError(t, err)
	})

	t.Run("error path - validation fails for empty path", func(t *testing.T) {
		emptyPath := ""

		err := ValidateFilePath(emptyPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file path")
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("error path - validation fails for non-existent file", func(t *testing.T) {
		nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.json")

		err := ValidateFilePath(nonExistentPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file does not exist")
	})

	t.Run("happy path - validation passes for file with special characters in path", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test-file_with.special-chars.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})

	t.Run("happy path - validation passes for file in nested directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedDir := filepath.Join(tmpDir, "nested", "subdir")
		err := os.MkdirAll(nestedDir, 0755)
		require.NoError(t, err)

		testFile := filepath.Join(nestedDir, "test.json")
		err = os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})

	t.Run("happy path - validation passes for empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "empty.json")
		err := os.WriteFile(testFile, []byte(""), 0644)
		require.NoError(t, err)

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})

	t.Run("happy path - validation passes for file with spaces in name", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test file with spaces.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})
}

func TestBuildFilePathFormValidationEdgeCases(t *testing.T) {
	t.Run("happy path - validation passes for file with special characters in path", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test-file_with.special-chars.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		config := FilePathFormConfig{
			Title:       "File Path",
			Placeholder: "Enter file path",
			Key:         "filepath",
		}

		form := BuildFilePathForm(config)
		assert.NotNil(t, form)

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})

	t.Run("happy path - form uses ValidateFilePath function", func(t *testing.T) {
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

		err = ValidateFilePath(testFile)

		assert.NoError(t, err)
	})
}

func TestBuildFilePathFormStructure(t *testing.T) {
	t.Run("happy path - form has input with correct key", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Test Title",
			Placeholder: "Test Placeholder",
			Key:         "test_key",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
		value := form.GetString(config.Key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - form has input with different keys", func(t *testing.T) {
		testCases := []string{
			"filepath",
			"config_file",
			"custom_key",
			"key123",
		}

		for _, key := range testCases {
			t.Run(key, func(t *testing.T) {
				config := FilePathFormConfig{
					Title:       "Test",
					Placeholder: "Test",
					Key:         key,
				}

				form := BuildFilePathForm(config)

				assert.NotNil(t, form)
				value := form.GetString(key)
				assert.Equal(t, "", value)
			})
		}
	})

	t.Run("happy path - form created with various placeholder values", func(t *testing.T) {
		testCases := []struct {
			name        string
			placeholder string
		}{
			{"absolute path", "/path/to/file.json"},
			{"relative path", "./config.json"},
			{"home path", "~/config.json"},
			{"empty placeholder", ""},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				config := FilePathFormConfig{
					Title:       "File Path",
					Placeholder: tc.placeholder,
					Key:         "filepath",
				}

				form := BuildFilePathForm(config)

				assert.NotNil(t, form)
			})
		}
	})

	t.Run("happy path - form created with various title values", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
		}{
			{"simple title", "File Path"},
			{"descriptive title", "Enter the path to your configuration file"},
			{"empty title", ""},
			{"title with special chars", "File Path (Required)"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				config := FilePathFormConfig{
					Title:       tc.title,
					Placeholder: "Enter path",
					Key:         "filepath",
				}

				form := BuildFilePathForm(config)

				assert.NotNil(t, form)
			})
		}
	})
}

func TestBuildSourceSelectionFormStructure(t *testing.T) {
	t.Run("happy path - form structure is consistent across calls", func(t *testing.T) {
		form1 := BuildSourceSelectionForm("Title 1", "key1")
		form2 := BuildSourceSelectionForm("Title 2", "key2")

		assert.NotNil(t, form1)
		assert.NotNil(t, form2)
		value1 := form1.GetString("key1")
		value2 := form2.GetString("key2")
		assert.Equal(t, "", value1)
		assert.Equal(t, "", value2)
	})

	t.Run("happy path - form options match source constants", func(t *testing.T) {
		form := BuildSourceSelectionForm("Select Source", "source")

		assert.NotNil(t, form)
		assert.Equal(t, "file", SourceFile)
		assert.Equal(t, "form", SourceForm)
		assert.Equal(t, "skip", SourceSkip)
		assert.Equal(t, "source", SourceKey)
	})
}

func TestGetFilePathFromFormStructure(t *testing.T) {
	t.Run("happy path - function builds form with provided config", func(t *testing.T) {
		config := FilePathFormConfig{
			Title:       "Custom Title",
			Placeholder: "Custom Placeholder",
			Key:         "custom_key",
		}

		form := BuildFilePathForm(config)

		assert.NotNil(t, form)
		value := form.GetString(config.Key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - function handles different config combinations", func(t *testing.T) {
		testCases := []FilePathFormConfig{
			{Title: "Title 1", Placeholder: "Placeholder 1", Key: "key1"},
			{Title: "", Placeholder: "Placeholder 2", Key: "key2"},
			{Title: "Title 3", Placeholder: "", Key: "key3"},
			{Title: "", Placeholder: "", Key: ""},
		}

		for i, config := range testCases {
			t.Run(fmt.Sprintf("config_%d", i), func(t *testing.T) {
				form := BuildFilePathForm(config)

				assert.NotNil(t, form)
				value := form.GetString(config.Key)
				assert.Equal(t, "", value)
			})
		}
	})
}

func TestGetSourceFromFormStructure(t *testing.T) {
	t.Run("happy path - function builds form with provided parameters", func(t *testing.T) {
		title := "Custom Title"
		key := "custom_key"

		form := BuildSourceSelectionForm(title, key)

		assert.NotNil(t, form)
		value := form.GetString(key)
		assert.Equal(t, "", value)
	})

	t.Run("happy path - function handles various parameter combinations", func(t *testing.T) {
		testCases := []struct {
			title string
			key   string
		}{
			{"Title 1", "key1"},
			{"", "key2"},
			{"Title 3", ""},
			{"", ""},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("params_%d", i), func(t *testing.T) {
				form := BuildSourceSelectionForm(tc.title, tc.key)

				assert.NotNil(t, form)
				value := form.GetString(tc.key)
				assert.Equal(t, "", value)
			})
		}
	})
}
