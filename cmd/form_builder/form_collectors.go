package form_builder

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/huh"

	"inzibat/config"
)

func CollectHeadersFromForm() (http.Header, error) {
	headers := make(http.Header)

	for {
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

		if err := headerForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect header: %w", err)
		}

		key := headerForm.GetString("key")
		value := headerForm.GetString("value")
		headers.Set(key, value)

		continueForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("continue").
					Title("Add another header?"),
			),
		)

		if err := continueForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueForm.GetBool("continue") {
			break
		}
	}

	return headers, nil
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

	if err := bodyForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to collect body field: %w", err)
	}

	body := make(config.HttpBody)
	key := bodyForm.GetString("key")
	value := bodyForm.GetString("value")
	body[key] = value

	for {
		continueForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("continue").
					Title("Add another body field?"),
			),
		)

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

	if err := bodyStringForm.Run(); err != nil {
		return "", fmt.Errorf("failed to get body string: %w", err)
	}

	return bodyStringForm.GetString("bodyString"), nil
}
