package form_builder

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStatusCode(t *testing.T) {
	t.Run("happy path - valid status code", func(t *testing.T) {
		validCodes := []string{"200", "201", "404", "500", "100", "511"}

		for _, code := range validCodes {
			err := ValidateStatusCode(code)
			assert.NoError(t, err, "code: %s", code)
		}
	})

	t.Run("error path - invalid status code format", func(t *testing.T) {
		invalidCodes := []string{"abc", "not-a-number", "", "12.5"}

		for _, code := range invalidCodes {
			err := ValidateStatusCode(code)
			assert.Error(t, err, "code: %s", code)
			assert.Contains(t, err.Error(), "invalid status code")
		}
	})

	t.Run("error path - status code out of range (too low)", func(t *testing.T) {
		code := "99"

		err := ValidateStatusCode(code)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status code must be between 100 and 511")
	})

	t.Run("error path - status code out of range (too high)", func(t *testing.T) {
		code := "512"

		err := ValidateStatusCode(code)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status code must be between 100 and 511")
	})

	t.Run("happy path - boundary values", func(t *testing.T) {
		minCode := "100"
		maxCode := "511"

		assert.NoError(t, ValidateStatusCode(minCode))
		assert.NoError(t, ValidateStatusCode(maxCode))
	})
}

func TestValidatePath(t *testing.T) {
	t.Run("happy path - valid path starting with /", func(t *testing.T) {
		validPaths := []string{"/", "/test", "/api/users", "/path/to/resource"}

		for _, path := range validPaths {
			err := ValidatePath(path)
			assert.NoError(t, err, "path: %s", path)
		}
	})

	t.Run("error path - path does not start with /", func(t *testing.T) {
		invalidPaths := []string{"test", "api/users", "path/to/resource", ""}

		for _, path := range invalidPaths {
			err := ValidatePath(path)
			assert.Error(t, err, "path: %s", path)
			assert.Contains(t, err.Error(), "route path must start with '/'")
		}
	})
}

func TestValidateHost(t *testing.T) {
	t.Run("happy path - valid URL", func(t *testing.T) {
		validHosts := []string{
			"http://localhost",
			"https://example.com",
			"http://localhost:8080",
			"https://api.example.com:443",
		}

		for _, host := range validHosts {
			err := ValidateHost(host)
			assert.NoError(t, err, "host: %s", host)
		}
	})

	t.Run("error path - invalid URL", func(t *testing.T) {
		invalidHosts := []string{
			"not a url",
			"://invalid",
		}

		for _, host := range invalidHosts {
			_, parseErr := url.Parse(host)
			if parseErr != nil {
				err := ValidateHost(host)
				assert.Error(t, err, "host: %s", host)
				assert.Contains(t, err.Error(), "invalid hostname")
			}
		}
	})

	t.Run("happy path - empty string (url.Parse allows it)", func(t *testing.T) {
		host := ""

		err := ValidateHost(host)

		_, parseErr := url.Parse(host)
		if parseErr == nil {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	})
}

func TestValidateNonEmpty(t *testing.T) {
	t.Run("happy path - non-empty value", func(t *testing.T) {
		nonEmptyValues := []string{"test", "value", "123", "a"}

		for _, value := range nonEmptyValues {
			err := ValidateNonEmpty(value, "field")
			assert.NoError(t, err, "value: %s", value)
		}
	})

	t.Run("error path - empty value", func(t *testing.T) {
		emptyValue := ""

		err := ValidateNonEmpty(emptyValue, "fieldName")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fieldName")
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("happy path - custom field name in error", func(t *testing.T) {
		emptyValue := ""
		fieldName := "customField"

		err := ValidateNonEmpty(emptyValue, fieldName)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fieldName)
	})
}
