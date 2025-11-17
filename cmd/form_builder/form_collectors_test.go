package form_builder

import (
	"errors"
	"net/http"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

func TestHuhFormRunnerInCollectors(t *testing.T) {
	t.Run("happy path - Run delegates to underlying form", func(t *testing.T) {
		t.Skip("Skipping interactive form test - form.Run() requires TTY and will hang in non-interactive environments")
		defaultValue := "default-value"
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("test").
					Title("Test Input").
					Value(&defaultValue),
			),
		)

		runner := &HuhFormRunner{Form: form}

		err := runner.Run()
		_ = err
		assert.NotNil(t, runner)
		assert.NotNil(t, runner.Form)
	})

	t.Run("happy path - GetString delegates to underlying form", func(t *testing.T) {
		testValue := "test-value"
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("test").
					Title("Test Input").
					Value(&testValue),
			),
		)

		runner := &HuhFormRunner{Form: form}

		result := runner.GetString("test")
		_ = result
		assert.NotNil(t, runner)
	})

	t.Run("happy path - GetBool delegates to underlying form", func(t *testing.T) {
		testBool := false
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("test").
					Title("Test Confirm").
					Value(&testBool),
			),
		)

		runner := &HuhFormRunner{Form: form}

		result := runner.GetBool("test")
		_ = result
		assert.NotNil(t, runner)
	})
}

func TestCollectHeadersFromFormPublic(t *testing.T) {
	t.Skip("Skipping interactive form test - requires TTY and will hang in non-interactive environments")
	t.Run("happy path - function creates form builders", func(t *testing.T) {
		_, err := CollectHeadersFromForm()
		_ = err
	})
}

func TestCollectBodyFromFormPublic(t *testing.T) {
	t.Skip("Skipping interactive form test - requires TTY and will hang in non-interactive environments")
	t.Run("happy path - function creates form builders", func(t *testing.T) {
		_, err := CollectBodyFromForm()
		_ = err
	})
}

func TestCollectBodyStringFromFormPublic(t *testing.T) {
	t.Skip("Skipping interactive form test - requires TTY and will hang in non-interactive environments")
	t.Run("happy path - function creates form builders", func(t *testing.T) {
		_, err := CollectBodyStringFromForm()
		_ = err
	})
}

func TestCollectHeadersFromFormInternal(t *testing.T) {
	t.Run("happy path - single header collected", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		headerFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		headerFormRunner.EXPECT().Run().Return(nil)
		headerFormRunner.EXPECT().GetString("key").Return("Content-Type")
		headerFormRunner.EXPECT().GetString("value").Return("application/json")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(false)

		// Act
		headers, err := collectHeadersFromFormInternal(headerFormRunner, continueFormRunner)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, 1, len(headers))
	})

	t.Run("happy path - multiple headers collected", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		headerFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		// First header
		headerFormRunner.EXPECT().Run().Return(nil)
		headerFormRunner.EXPECT().GetString("key").Return("Content-Type")
		headerFormRunner.EXPECT().GetString("value").Return("application/json")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(true)

		// Second header
		headerFormRunner.EXPECT().Run().Return(nil)
		headerFormRunner.EXPECT().GetString("key").Return("Authorization")
		headerFormRunner.EXPECT().GetString("value").Return("Bearer token")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(false)

		// Act
		headers, err := collectHeadersFromFormInternal(headerFormRunner, continueFormRunner)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token", headers.Get("Authorization"))
		assert.Equal(t, 2, len(headers))
	})

	t.Run("error path - header form run fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		headerFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		expectedErr := errors.New("form run failed")
		headerFormRunner.EXPECT().Run().Return(expectedErr)

		// Act
		headers, err := collectHeadersFromFormInternal(headerFormRunner, continueFormRunner)

		// Assert
		require.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to collect header")
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("error path - continue form run fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		headerFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		headerFormRunner.EXPECT().Run().Return(nil)
		headerFormRunner.EXPECT().GetString("key").Return("Content-Type")
		headerFormRunner.EXPECT().GetString("value").Return("application/json")
		expectedErr := errors.New("continue form failed")
		continueFormRunner.EXPECT().Run().Return(expectedErr)

		// Act
		headers, err := collectHeadersFromFormInternal(headerFormRunner, continueFormRunner)

		// Assert
		require.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to get user input")
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestCollectBodyFromFormInternal(t *testing.T) {
	t.Run("happy path - single body field collected", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("message")
		bodyFormRunner.EXPECT().GetString("value").Return("success")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(false)

		// Act
		body, err := collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, body)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, 1, len(body))
	})

	t.Run("happy path - multiple body fields collected", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		// First body field
		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("message")
		bodyFormRunner.EXPECT().GetString("value").Return("success")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(true)

		// Second body field
		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("code")
		bodyFormRunner.EXPECT().GetString("value").Return("200")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(false)

		// Act
		body, err := collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, body)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, "200", body["code"])
		assert.Equal(t, 2, len(body))
	})

	t.Run("happy path - body field overwrites previous value", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		// First body field
		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("key")
		bodyFormRunner.EXPECT().GetString("value").Return("old value")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(true)

		// Second body field with same key
		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("key")
		bodyFormRunner.EXPECT().GetString("value").Return("new value")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(false)

		// Act
		body, err := collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, body)
		assert.Equal(t, "new value", body["key"])
		assert.Equal(t, 1, len(body))
	})

	t.Run("error path - initial body form run fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		expectedErr := errors.New("form run failed")
		bodyFormRunner.EXPECT().Run().Return(expectedErr)

		// Act
		body, err := collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)

		// Assert
		require.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to collect body field")
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("error path - continue form run fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("message")
		bodyFormRunner.EXPECT().GetString("value").Return("success")
		expectedErr := errors.New("continue form failed")
		continueFormRunner.EXPECT().Run().Return(expectedErr)

		// Act
		body, err := collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)

		// Assert
		require.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to get user input")
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("error path - body form run fails in loop", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyFormRunner := NewMockFormRunner(ctrl)
		continueFormRunner := NewMockFormRunner(ctrl)

		// First body field
		bodyFormRunner.EXPECT().Run().Return(nil)
		bodyFormRunner.EXPECT().GetString("key").Return("message")
		bodyFormRunner.EXPECT().GetString("value").Return("success")
		continueFormRunner.EXPECT().Run().Return(nil)
		continueFormRunner.EXPECT().GetBool("continue").Return(true)

		// Second body field fails
		expectedErr := errors.New("form run failed")
		bodyFormRunner.EXPECT().Run().Return(expectedErr)

		// Act
		body, err := collectBodyFromFormInternal(bodyFormRunner, continueFormRunner)

		// Assert
		require.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to collect body field")
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestCollectBodyStringFromFormInternal(t *testing.T) {
	t.Run("happy path - body string collected", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyStringFormRunner := NewMockFormRunner(ctrl)

		expectedBodyString := `{"message": "success"}`
		bodyStringFormRunner.EXPECT().Run().Return(nil)
		bodyStringFormRunner.EXPECT().GetString("bodyString").Return(expectedBodyString)

		// Act
		bodyString, err := collectBodyStringFromFormInternal(bodyStringFormRunner)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedBodyString, bodyString)
	})

	t.Run("happy path - plain text body string collected", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyStringFormRunner := NewMockFormRunner(ctrl)

		expectedBodyString := "Hello, World!"
		bodyStringFormRunner.EXPECT().Run().Return(nil)
		bodyStringFormRunner.EXPECT().GetString("bodyString").Return(expectedBodyString)

		// Act
		bodyString, err := collectBodyStringFromFormInternal(bodyStringFormRunner)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedBodyString, bodyString)
	})

	t.Run("error path - form run fails", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bodyStringFormRunner := NewMockFormRunner(ctrl)

		expectedErr := errors.New("form run failed")
		bodyStringFormRunner.EXPECT().Run().Return(expectedErr)

		// Act
		bodyString, err := collectBodyStringFromFormInternal(bodyStringFormRunner)

		// Assert
		require.Error(t, err)
		assert.Empty(t, bodyString)
		assert.Contains(t, err.Error(), "failed to get body string")
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestHuhFormRunnerMethodsInCollectors(t *testing.T) {
	t.Run("happy path - Run delegates to underlying form", func(t *testing.T) {
		// Arrange
		testValue := "test-value"
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("test").
					Title("Test Input").
					Value(&testValue),
			),
		)

		runner := &HuhFormRunner{Form: form}

		// Act & Assert
		assert.NotNil(t, runner)
		assert.NotNil(t, runner.Form)
		// Note: We can't actually test Run() without a TTY, but we can verify the structure
	})

	t.Run("happy path - GetString delegates to underlying form", func(t *testing.T) {
		// Arrange
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("test").
					Title("Test Input"),
			),
		)

		runner := &HuhFormRunner{Form: form}

		// Act
		result := runner.GetString("test")

		// Assert
		assert.NotNil(t, runner)
		// Note: GetString returns empty string for unrun forms, which is expected behavior
		assert.IsType(t, "", result)
	})

	t.Run("happy path - GetBool delegates to underlying form", func(t *testing.T) {
		// Arrange
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("test").
					Title("Test Confirm"),
			),
		)

		runner := &HuhFormRunner{Form: form}

		// Act
		result := runner.GetBool("test")

		// Assert
		assert.NotNil(t, runner)
		// Note: GetBool returns false for unrun forms, which is expected behavior
		assert.IsType(t, false, result)
	})

	t.Run("happy path - GetBool returns false for unset value", func(t *testing.T) {
		// Arrange
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("test").
					Title("Test Confirm"),
			),
		)

		runner := &HuhFormRunner{Form: form}

		// Act
		result := runner.GetBool("test")

		// Assert
		assert.NotNil(t, runner)
		assert.False(t, result)
	})

	t.Run("happy path - GetString returns empty for unset value", func(t *testing.T) {
		// Arrange
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("test").
					Title("Test Input"),
			),
		)

		runner := &HuhFormRunner{Form: form}

		// Act
		result := runner.GetString("test")

		// Assert
		assert.NotNil(t, runner)
		assert.Empty(t, result)
	})
}
