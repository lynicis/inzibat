package form_builder

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/huh"

	"inzibat/config"
)

type FormRunner interface {
	Run() error
	GetString(key string) string
	GetBool(key string) bool
}

type huhFormRunner struct {
	form *huh.Form
}

func (r *huhFormRunner) Run() error {
	return r.form.Run()
}

func (r *huhFormRunner) GetString(key string) string {
	return r.form.GetString(key)
}

func (r *huhFormRunner) GetBool(key string) bool {
	return r.form.GetBool(key)
}

func collectHeadersFromFormWithRunner(
	buildHeaderForm func() FormRunner,
	buildContinueForm func() FormRunner,
) (http.Header, error) {
	headers := make(http.Header)

	for {
		headerForm := buildHeaderForm()

		if err := headerForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect header: %w", err)
		}

		key := headerForm.GetString("key")
		value := headerForm.GetString("value")
		headers.Set(key, value)

		continueForm := buildContinueForm()

		if err := continueForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueForm.GetBool("continue") {
			break
		}
	}

	return headers, nil
}

func CollectHeadersFromForm() (http.Header, error) {
	return collectHeadersFromFormWithRunner(
		func() FormRunner {
			return &huhFormRunner{
				form: huh.NewForm(
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
				),
			}
		},
		func() FormRunner {
			return &huhFormRunner{
				form: huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Key("continue").
							Title("Add another header?"),
					),
				),
			}
		},
	)
}

func collectBodyFromFormWithRunner(
	buildBodyForm func() FormRunner,
	buildContinueForm func() FormRunner,
) (config.HttpBody, error) {
	bodyForm := buildBodyForm()

	if err := bodyForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to collect body field: %w", err)
	}

	body := make(config.HttpBody)
	key := bodyForm.GetString("key")
	value := bodyForm.GetString("value")
	body[key] = value

	for {
		continueForm := buildContinueForm()

		if err := continueForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueForm.GetBool("continue") {
			break
		}

		if err := bodyForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect body field: %w", err)
		}

		key := bodyForm.GetString("key")
		value := bodyForm.GetString("value")
		body[key] = value
	}

	return body, nil
}

func CollectBodyFromForm() (config.HttpBody, error) {
	return collectBodyFromFormWithRunner(
		func() FormRunner {
			return &huhFormRunner{
				form: huh.NewForm(
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
				),
			}
		},
		func() FormRunner {
			return &huhFormRunner{
				form: huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Key("continue").
							Title("Add another body field?"),
					),
				),
			}
		},
	)
}

func collectBodyStringFromFormWithRunner(buildBodyStringForm func() FormRunner) (string, error) {
	bodyStringForm := buildBodyStringForm()

	if err := bodyStringForm.Run(); err != nil {
		return "", fmt.Errorf("failed to get body string: %w", err)
	}

	return bodyStringForm.GetString("bodyString"), nil
}

func CollectBodyStringFromForm() (string, error) {
	return collectBodyStringFromFormWithRunner(
		func() FormRunner {
			return &huhFormRunner{
				form: huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Key("bodyString").
							Title("Body String").
							Placeholder(`{"message": "success"}`).
							Validate(func(s string) error {
								return ValidateNonEmpty(s, "body string")
							}),
					),
				),
			}
		},
	)
}
