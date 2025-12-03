package form_builder

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/lynicis/inzibat/config"
)

type FormRunner interface {
	Run() error
	GetString(key string) string
	GetBool(key string) bool
}

// HuhFormRunner is a sugar wrapper that makes logic testable
type HuhFormRunner struct {
	Form *huh.Form
}

func (r *HuhFormRunner) Run() error {
	return r.Form.Run()
}

func (r *HuhFormRunner) GetString(key string) string {
	return r.Form.GetString(key)
}

func (r *HuhFormRunner) GetBool(key string) bool {
	return r.Form.GetBool(key)
}

func collectHeadersFromFormInternal(
	headerFormCreator func() *huh.Form,
	continueFormCreator func() *huh.Form,
) (http.Header, error) {
	headers := make(http.Header)

	for {
		// Create a fresh form instance for each iteration
		headerForm := headerFormCreator()
		if err := headerForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect header: %w", err)
		}

		key := headerForm.GetString("key")
		value := headerForm.GetString("value")
		headers.Set(key, value)

		// Create a fresh continue form for each iteration
		continueForm := continueFormCreator()
		if err := continueForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueForm.GetBool("continue") {
			break
		}
	}

	return headers, nil
}

func createHeaderForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("key").
				Title("Header Key").
				Placeholder("Content-Type").
				Validate(func(s string) error {
					return ValidateNonEmpty(s, "header key")
				}),
			huh.NewInput().
				Key("value").
				Title("Header Value").
				Placeholder("application/json").
				Validate(func(s string) error {
					return ValidateNonEmpty(s, "header value")
				}),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
}

func createContinueForm(message string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("continue").
				Title(message).
				Affirmative("Yes").
				Negative("No"),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
}

func CollectHeadersFromForm() (http.Header, error) {
	return collectHeadersFromFormInternal(
		createHeaderForm,
		func() *huh.Form { return createContinueForm("Add another header?") },
	)
}

func collectBodyFromFormInternal(
	bodyFormCreator func() *huh.Form,
	continueFormCreator func() *huh.Form,
) (config.HttpBody, error) {
	body := make(config.HttpBody)

	// Create first body form
	bodyForm := bodyFormCreator()
	if err := bodyForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to collect body field: %w", err)
	}

	key := bodyForm.GetString("key")
	value := bodyForm.GetString("value")
	body[key] = value

	for {
		// Create a fresh continue form for each iteration
		continueForm := continueFormCreator()
		if err := continueForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueForm.GetBool("continue") {
			break
		}

		// Create a fresh body form for each iteration
		bodyForm := bodyFormCreator()
		if err := bodyForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect body field: %w", err)
		}

		key := bodyForm.GetString("key")
		value := bodyForm.GetString("value")
		body[key] = value
	}

	return body, nil
}

func createBodyForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("key").
				Title("Body Key").
				Placeholder("message").
				Validate(func(s string) error {
					return ValidateNonEmpty(s, "body key")
				}),
			huh.NewInput().
				Key("value").
				Title("Body Value").
				Placeholder("success").
				Validate(func(s string) error {
					return ValidateNonEmpty(s, "body value")
				}),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
}

func CollectBodyFromForm() (config.HttpBody, error) {
	return collectBodyFromFormInternal(
		createBodyForm,
		func() *huh.Form { return createContinueForm("Add another body field?") },
	)
}

func collectBodyStringFromFormInternal(bodyStringFormRunner FormRunner) (string, error) {
	if err := bodyStringFormRunner.Run(); err != nil {
		return "", fmt.Errorf("failed to get body string: %w", err)
	}

	return bodyStringFormRunner.GetString("bodyString"), nil
}

func CollectBodyStringFromForm() (string, error) {
	bodyStringForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("bodyString").
				Title("Body String").
				Placeholder(`{"message": "success"}`).
				Validate(func(s string) error {
					return ValidateNonEmpty(s, "body string")
				}),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

	bodyStringFormRunner := &HuhFormRunner{Form: bodyStringForm}

	return collectBodyStringFromFormInternal(bodyStringFormRunner)
}
