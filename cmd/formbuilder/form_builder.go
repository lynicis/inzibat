package formbuilder

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type FilePathFormConfig struct {
	Title       string
	Placeholder string
	Key         string
}

// NewHuhForm applies the terminal wiring every interactive form shares.
func NewHuhForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
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

	return NewHuhForm(huh.NewGroup(input))
}

func BuildSourceSelectionForm(title string, key string) *huh.Form {
	return NewHuhForm(
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
