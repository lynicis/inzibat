package form_builder

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

type SourceOption struct {
	Title string
	Value string
}

type FilePathFormConfig struct {
	Title       string
	Placeholder string
	Key         string
}

type FormBuilder struct {
	title       string
	key         string
	placeholder string
	validateFn  func(string) error
}

func NewFormBuilder() *FormBuilder {
	return &FormBuilder{}
}

func (b *FormBuilder) WithTitle(title string) *FormBuilder {
	b.title = title
	return b
}

func (b *FormBuilder) WithKey(key string) *FormBuilder {
	b.key = key
	return b
}

func (b *FormBuilder) WithPlaceholder(placeholder string) *FormBuilder {
	b.placeholder = placeholder
	return b
}

func (b *FormBuilder) WithValidation(validateFn func(string) error) *FormBuilder {
	b.validateFn = validateFn
	return b
}

func (b *FormBuilder) BuildInputForm() *huh.Form {
	input := huh.NewInput().
		Key(b.key).
		Title(b.title).
		Placeholder(b.placeholder)

	if b.validateFn != nil {
		input = input.Validate(b.validateFn)
	}

	return huh.NewForm(huh.NewGroup(input))
}

func BuildFilePathForm(config FilePathFormConfig) *huh.Form {
	return NewFormBuilder().
		WithKey(config.Key).
		WithTitle(config.Title).
		WithPlaceholder(config.Placeholder).
		WithValidation(func(s string) error {
			if err := ValidateNonEmpty(s, "file path"); err != nil {
				return err
			}
			if _, err := os.Stat(s); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist")
			}
			return nil
		}).
		BuildInputForm()
}

func BuildSourceSelectionForm(title string, key string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key(key).
				Title(title).
				Options([]huh.Option[string]{
					{Key: "File", Value: SourceFile},
					{Key: "Form", Value: SourceForm},
					{Key: "Skip", Value: SourceSkip},
				}...),
		),
	)
}

func GetFilePathFromForm(config FilePathFormConfig) (string, error) {
	form := BuildFilePathForm(config)
	if err := form.Run(); err != nil {
		return "", fmt.Errorf("failed to get file path: %w", err)
	}
	return form.GetString(config.Key), nil
}

func GetSourceFromForm(title string, key string) (string, error) {
	form := BuildSourceSelectionForm(title, key)
	if err := form.Run(); err != nil {
		return "", fmt.Errorf("failed to select source: %w", err)
	}
	return form.GetString(key), nil
}
