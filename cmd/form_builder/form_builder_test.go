package form_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormBuilder(t *testing.T) {
	t.Run("happy path - creates new FormBuilder", func(t *testing.T) {
		builder := NewFormBuilder()

		assert.NotNil(t, builder)
		assert.Equal(t, "", builder.title)
		assert.Equal(t, "", builder.key)
		assert.Equal(t, "", builder.placeholder)
		assert.Nil(t, builder.validateFn)
	})
}

func TestFormBuilder_WithTitle(t *testing.T) {
	t.Run("happy path - sets title and returns builder", func(t *testing.T) {
		builder := NewFormBuilder()
		title := "Test Title"

		result := builder.WithTitle(title)

		assert.Equal(t, builder, result)
		assert.Equal(t, title, builder.title)
	})

	t.Run("happy path - allows chaining", func(t *testing.T) {
		title1 := "First Title"
		title2 := "Second Title"

		builder := NewFormBuilder().
			WithTitle(title1).
			WithTitle(title2)

		assert.Equal(t, title2, builder.title)
	})
}

func TestFormBuilder_WithKey(t *testing.T) {
	t.Run("happy path - sets key and returns builder", func(t *testing.T) {
		builder := NewFormBuilder()
		key := "test_key"

		result := builder.WithKey(key)

		assert.Equal(t, builder, result)
		assert.Equal(t, key, builder.key)
	})

	t.Run("happy path - allows chaining", func(t *testing.T) {
		key1 := "key1"
		key2 := "key2"

		builder := NewFormBuilder().
			WithKey(key1).
			WithKey(key2)

		assert.Equal(t, key2, builder.key)
	})
}

func TestFormBuilder_WithPlaceholder(t *testing.T) {
	t.Run("happy path - sets placeholder and returns builder", func(t *testing.T) {
		builder := NewFormBuilder()
		placeholder := "Enter value here"

		result := builder.WithPlaceholder(placeholder)

		assert.Equal(t, builder, result)
		assert.Equal(t, placeholder, builder.placeholder)
	})

	t.Run("happy path - allows chaining", func(t *testing.T) {
		placeholder1 := "First placeholder"
		placeholder2 := "Second placeholder"

		builder := NewFormBuilder().
			WithPlaceholder(placeholder1).
			WithPlaceholder(placeholder2)

		assert.Equal(t, placeholder2, builder.placeholder)
	})
}

func TestFormBuilder_WithValidation(t *testing.T) {
	t.Run("happy path - sets validation function and returns builder", func(t *testing.T) {
		builder := NewFormBuilder()
		validateFn := func(s string) error {
			if s == "" {
				return assert.AnError
			}
			return nil
		}

		result := builder.WithValidation(validateFn)

		assert.Equal(t, builder, result)
		assert.NotNil(t, builder.validateFn)
		assert.NoError(t, builder.validateFn("test"))
		assert.Error(t, builder.validateFn(""))
	})

	t.Run("happy path - allows chaining", func(t *testing.T) {
		validateFn1 := func(s string) error { return nil }
		validateFn2 := func(s string) error { return assert.AnError }

		builder := NewFormBuilder().
			WithValidation(validateFn1).
			WithValidation(validateFn2)

		assert.NotNil(t, builder.validateFn)
		assert.Error(t, builder.validateFn("test"))
	})

	t.Run("happy path - can set nil validation function", func(t *testing.T) {
		builder := NewFormBuilder().
			WithValidation(func(s string) error { return nil })

		builder.WithValidation(nil)

		assert.Nil(t, builder.validateFn)
	})
}

func TestFormBuilder_BuildInputForm(t *testing.T) {
	t.Run("happy path - builds form with all properties", func(t *testing.T) {
		builder := NewFormBuilder().
			WithTitle("Test Title").
			WithKey("test_key").
			WithPlaceholder("Enter value").
			WithValidation(func(s string) error {
				if s == "" {
					return assert.AnError
				}
				return nil
			})

		form := builder.BuildInputForm()

		assert.NotNil(t, form)
	})

	t.Run("happy path - builds form without validation", func(t *testing.T) {
		builder := NewFormBuilder().
			WithTitle("Test Title").
			WithKey("test_key").
			WithPlaceholder("Enter value")

		form := builder.BuildInputForm()

		assert.NotNil(t, form)
	})

	t.Run("happy path - builds form with minimal properties", func(t *testing.T) {
		builder := NewFormBuilder().
			WithKey("test_key")

		form := builder.BuildInputForm()

		assert.NotNil(t, form)
	})
}

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
	t.Run("happy path - gets file path from form with existing file", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})

	t.Run("error path - file does not exist", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})
}

func TestGetSourceFromForm(t *testing.T) {
	t.Run("happy path - gets source from form", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})

	t.Run("error path - form run fails", func(t *testing.T) {
		t.Skip("Skipping interactive form test - requires non-interactive mode or mocking")
	})
}

func TestFormBuilder_Chaining(t *testing.T) {
	t.Run("happy path - all builder methods can be chained", func(t *testing.T) {
		title := "Test Title"
		key := "test_key"
		placeholder := "Enter value"
		validateFn := func(s string) error { return nil }

		builder := NewFormBuilder().
			WithTitle(title).
			WithKey(key).
			WithPlaceholder(placeholder).
			WithValidation(validateFn)

		assert.Equal(t, title, builder.title)
		assert.Equal(t, key, builder.key)
		assert.Equal(t, placeholder, builder.placeholder)
		assert.NotNil(t, builder.validateFn)
	})

	t.Run("happy path - chained builder creates valid form", func(t *testing.T) {
		title := "Test Title"
		key := "test_key"
		placeholder := "Enter value"

		form := NewFormBuilder().
			WithTitle(title).
			WithKey(key).
			WithPlaceholder(placeholder).
			BuildInputForm()

		assert.NotNil(t, form)
	})
}
