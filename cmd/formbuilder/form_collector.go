package formbuilder

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

type keyValuePair struct {
	Key   string
	Value string
}

func collectPairs(
	entryCreator func() FormRunner,
	continueCreator func() FormRunner,
	entryTitle string,
) ([]keyValuePair, error) {
	var pairs []keyValuePair

	for {
		entryForm := entryCreator()
		if err := entryForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to collect %s: %w", entryTitle, err)
		}

		pairs = append(pairs, keyValuePair{
			Key:   entryForm.GetString("key"),
			Value: entryForm.GetString("value"),
		})

		continueForm := continueCreator()
		if err := continueForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get user input: %w", err)
		}

		if !continueForm.GetBool("continue") {
			break
		}
	}

	return pairs, nil
}

func collectHeadersFromFormInternal(
	headerFormCreator func() FormRunner,
	continueFormCreator func() FormRunner,
) (http.Header, error) {
	pairs, err := collectPairs(headerFormCreator, continueFormCreator, "header")
	if err != nil {
		return nil, err
	}

	headers := make(http.Header, len(pairs))
	for _, pair := range pairs {
		headers.Set(pair.Key, pair.Value)
	}

	return headers, nil
}

func createHeaderForm() *huh.Form {
	return NewHuhForm(
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
}

func createContinueForm(message string) *huh.Form {
	return NewHuhForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("continue").
				Title(message).
				Affirmative("Yes").
				Negative("No"),
		),
	)
}

func CollectHeadersFromForm() (http.Header, error) {
	return collectHeadersFromFormInternal(
		func() FormRunner { return &HuhFormRunner{Form: createHeaderForm()} },
		func() FormRunner {
			return &HuhFormRunner{Form: createContinueForm("Add another header?")}
		},
	)
}

func collectBodyFromFormInternal(
	bodyFormCreator func() FormRunner,
	continueFormCreator func() FormRunner,
) (config.HttpBody, error) {
	pairs, err := collectPairs(bodyFormCreator, continueFormCreator, "body field")
	if err != nil {
		return nil, err
	}

	body := make(config.HttpBody, len(pairs))
	for _, pair := range pairs {
		body[pair.Key] = pair.Value
	}

	return body, nil
}

func createBodyForm() *huh.Form {
	return NewHuhForm(
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
}

func CollectBodyFromForm() (config.HttpBody, error) {
	return collectBodyFromFormInternal(
		func() FormRunner { return &HuhFormRunner{Form: createBodyForm()} },
		func() FormRunner {
			return &HuhFormRunner{Form: createContinueForm("Add another body field?")}
		},
	)
}

func collectBodyStringFromFormInternal(bodyStringFormRunner FormRunner) (string, error) {
	if err := bodyStringFormRunner.Run(); err != nil {
		return "", fmt.Errorf("failed to get body string: %w", err)
	}

	return bodyStringFormRunner.GetString("bodyString"), nil
}

func CollectBodyStringFromForm() (string, error) {
	bodyStringForm := NewHuhForm(
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
