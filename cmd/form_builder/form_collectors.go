package form_builder

import (
	"fmt"
	"net/http"

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
	headerFormRunner,
	continueFormRunner FormRunner,
) (http.Header, error) {
	headers := make(http.Header)

	for {
		if err := headerFormRunner.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect header: %w", err)
		}

		key := headerFormRunner.GetString("key")
		value := headerFormRunner.GetString("value")
		headers.Set(key, value)

		if err := continueFormRunner.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueFormRunner.GetBool("continue") {
			break
		}
	}

	return headers, nil
}

func CollectHeadersFromForm() (http.Header, error) {
	headerForm := huh.NewForm(
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
	)

	continueForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("continue").
				Title("Add another header?"),
		),
	)

	headerFormRunner := &HuhFormRunner{Form: headerForm}
	continueFormRunner := &HuhFormRunner{Form: continueForm}

	return collectHeadersFromFormInternal(headerFormRunner, continueFormRunner)
}

func collectBodyFromFormInternal(
	bodyFormRunner,
	continueFormRunner FormRunner,
) (config.HttpBody, error) {
	body := make(config.HttpBody)

	if err := bodyFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to collect body field: %w", err)
	}

	key := bodyFormRunner.GetString("key")
	value := bodyFormRunner.GetString("value")
	body[key] = value

	for {
		if err := continueFormRunner.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueFormRunner.GetBool("continue") {
			break
		}

		if err := bodyFormRunner.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect body field: %w", err)
		}

		key := bodyFormRunner.GetString("key")
		value := bodyFormRunner.GetString("value")
		body[key] = value
	}

	return body, nil
}

func CollectBodyFromForm() (config.HttpBody, error) {
	bodyForm := huh.NewForm(
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
	)

	continueForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("continue").
				Title("Add another body field?"),
		),
	)

	bodyFormRunner := &HuhFormRunner{Form: bodyForm}
	continueFormRunner := &HuhFormRunner{Form: continueForm}

	return collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)
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
	)

	bodyStringFormRunner := &HuhFormRunner{Form: bodyStringForm}

	return collectBodyStringFromFormInternal(bodyStringFormRunner)
}
