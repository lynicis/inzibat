package form_builder

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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

func ValidateFilePath(s string) error {
	if err := ValidateNonEmpty(s, "file path"); err != nil {
		return err
	}
	if _, err := os.Stat(s); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}
	return nil
}

func BuildFilePathForm(config FilePathFormConfig) *huh.Form {
	input := huh.NewInput().
		Key(config.Key).
		Title(config.Title).
		Placeholder(config.Placeholder).
		Validate(ValidateFilePath)

	return huh.NewForm(huh.NewGroup(input)).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
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
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
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
