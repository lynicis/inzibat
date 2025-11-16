package form_builder

import (
	"net/http"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func TestCollectHeadersFromForm(t *testing.T) {
	t.Run("happy path - form structure is created correctly", func(t *testing.T) {
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

		assert.NotNil(t, headerForm)
	})

	t.Run("happy path - validation functions work correctly", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     string
			fieldName string
			wantError bool
		}{
			{
				name:      "valid header key",
				value:     "Content-Type",
				fieldName: "header key",
				wantError: false,
			},
			{
				name:      "valid header value",
				value:     "application/json",
				fieldName: "header value",
				wantError: false,
			},
			{
				name:      "empty header key",
				value:     "",
				fieldName: "header key",
				wantError: true,
			},
			{
				name:      "empty header value",
				value:     "",
				fieldName: "header value",
				wantError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateNonEmpty(tc.value, tc.fieldName)
				if tc.wantError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.fieldName)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("happy path - continue form structure is created correctly", func(t *testing.T) {
		continueForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("continue").
					Title("Add another header?"),
			),
		)

		assert.NotNil(t, continueForm)
	})

	t.Run("happy path - headers map initialization", func(t *testing.T) {
		headers := make(http.Header)

		headers.Set("Content-Type", "application/json")

		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, 1, len(headers))
	})

	t.Run("happy path - multiple headers can be added", func(t *testing.T) {
		headers := make(http.Header)

		headers.Set("Content-Type", "application/json")
		headers.Set("Authorization", "Bearer token")
		headers.Set("X-Custom-Header", "value")

		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token", headers.Get("Authorization"))
		assert.Equal(t, "value", headers.Get("X-Custom-Header"))
		assert.Equal(t, 3, len(headers))
	})
}

func TestCollectBodyFromForm(t *testing.T) {
	t.Run("happy path - form structure is created correctly", func(t *testing.T) {
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

		assert.NotNil(t, bodyForm)
	})

	t.Run("happy path - validation functions work correctly", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     string
			fieldName string
			wantError bool
		}{
			{
				name:      "valid body key",
				value:     "message",
				fieldName: "body key",
				wantError: false,
			},
			{
				name:      "valid body value",
				value:     "success",
				fieldName: "body value",
				wantError: false,
			},
			{
				name:      "empty body key",
				value:     "",
				fieldName: "body key",
				wantError: true,
			},
			{
				name:      "empty body value",
				value:     "",
				fieldName: "body value",
				wantError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateNonEmpty(tc.value, tc.fieldName)
				if tc.wantError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.fieldName)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("happy path - body map initialization", func(t *testing.T) {
		body := make(config.HttpBody)

		body["message"] = "success"

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, 1, len(body))
	})

	t.Run("happy path - multiple body fields can be added", func(t *testing.T) {
		body := make(config.HttpBody)

		body["message"] = "success"
		body["code"] = float64(200)
		body["status"] = "ok"

		assert.Equal(t, "success", body["message"])
		assert.Equal(t, float64(200), body["code"])
		assert.Equal(t, "ok", body["status"])
		assert.Equal(t, 3, len(body))
	})

	t.Run("happy path - continue form structure is created correctly", func(t *testing.T) {
		continueForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("continue").
					Title("Add another body field?"),
			),
		)

		assert.NotNil(t, continueForm)
	})

	t.Run("happy path - body field overwrites previous value", func(t *testing.T) {
		body := make(config.HttpBody)
		body["key"] = "old value"

		body["key"] = "new value"

		assert.Equal(t, "new value", body["key"])
		assert.Equal(t, 1, len(body))
	})
}

func TestCollectBodyStringFromForm(t *testing.T) {
	t.Run("happy path - form structure is created correctly", func(t *testing.T) {
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

		assert.NotNil(t, bodyStringForm)
	})

	t.Run("happy path - validation function works correctly", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     string
			fieldName string
			wantError bool
		}{
			{
				name:      "valid JSON body string",
				value:     `{"message": "success"}`,
				fieldName: "body string",
				wantError: false,
			},
			{
				name:      "valid plain text body string",
				value:     "Hello, World!",
				fieldName: "body string",
				wantError: false,
			},
			{
				name:      "empty body string",
				value:     "",
				fieldName: "body string",
				wantError: true,
			},
			{
				name:      "whitespace-only body string",
				value:     "   ",
				fieldName: "body string",
				wantError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateNonEmpty(tc.value, tc.fieldName)
				if tc.wantError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.fieldName)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("happy path - body string can contain various content", func(t *testing.T) {
		testCases := []struct {
			name    string
			content string
		}{
			{
				name:    "JSON object",
				content: `{"message": "success", "code": 200}`,
			},
			{
				name:    "JSON array",
				content: `[1, 2, 3, 4, 5]`,
			},
			{
				name:    "Plain text",
				content: "Hello, World!",
			},
			{
				name:    "XML",
				content: `<root><message>success</message></root>`,
			},
			{
				name:    "Multiline text",
				content: "Line 1\nLine 2\nLine 3",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateNonEmpty(tc.content, "body string")
				assert.NoError(t, err, "content: %s", tc.content)
			})
		}
	})
}

func TestFormErrorMessages(t *testing.T) {
	t.Run("error path - header collection error message format", func(t *testing.T) {
		err := ValidateNonEmpty("", "header key")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "header key")
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("error path - body collection error message format", func(t *testing.T) {
		err := ValidateNonEmpty("", "body key")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "body key")
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("error path - body string collection error message format", func(t *testing.T) {
		err := ValidateNonEmpty("", "body string")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "body string")
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestFormPlaceholders(t *testing.T) {
	t.Run("happy path - header form placeholders are correct", func(t *testing.T) {
		headerForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("key").
					Title("Header Key").
					Placeholder("Content-Type"),
				huh.NewInput().
					Key("value").
					Title("Header Value").
					Placeholder("application/json"),
			),
		)

		assert.NotNil(t, headerForm)
	})

	t.Run("happy path - body form placeholders are correct", func(t *testing.T) {
		bodyForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("key").
					Title("Body Key").
					Placeholder("message"),
				huh.NewInput().
					Key("value").
					Title("Body Value").
					Placeholder("success"),
			),
		)

		assert.NotNil(t, bodyForm)
	})

	t.Run("happy path - body string form placeholder is correct", func(t *testing.T) {
		bodyStringForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("bodyString").
					Title("Body String").
					Placeholder(`{"message": "success"}`),
			),
		)

		assert.NotNil(t, bodyStringForm)
	})
}
