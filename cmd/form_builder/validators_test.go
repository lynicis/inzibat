package form_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStatusCode(t *testing.T) {
	t.Run("happy path - valid status codes", func(t *testing.T) {
		testCases := []struct {
			name     string
			codeStr  string
			expected error
		}{
			{
				name:     "status 100",
				codeStr:  "100",
				expected: nil,
			},
			{
				name:     "status 200",
				codeStr:  "200",
				expected: nil,
			},
			{
				name:     "status 404",
				codeStr:  "404",
				expected: nil,
			},
			{
				name:     "status 500",
				codeStr:  "500",
				expected: nil,
			},
			{
				name:     "status 511",
				codeStr:  "511",
				expected: nil,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateStatusCode(tc.codeStr)

				assert.NoError(t, err)
			})
		}
	})

	t.Run("error path - invalid status code format", func(t *testing.T) {
		testCases := []struct {
			name     string
			codeStr  string
			expected string
		}{
			{
				name:     "non-numeric string",
				codeStr:  "abc",
				expected: "invalid status code",
			},
			{
				name:     "empty string",
				codeStr:  "",
				expected: "invalid status code",
			},
			{
				name:     "mixed alphanumeric",
				codeStr:  "200abc",
				expected: "invalid status code",
			},
			{
				name:     "negative number as string",
				codeStr:  "-100",
				expected: "status code must be between 100 and 511",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateStatusCode(tc.codeStr)

				assert.Error(t, err)
				assert.Equal(t, tc.expected, err.Error())
			})
		}
	})

	t.Run("error path - status code out of range", func(t *testing.T) {
		testCases := []struct {
			name     string
			codeStr  string
			expected string
		}{
			{
				name:     "status code too low",
				codeStr:  "99",
				expected: "status code must be between 100 and 511",
			},
			{
				name:     "status code too high",
				codeStr:  "512",
				expected: "status code must be between 100 and 511",
			},
			{
				name:     "status code zero",
				codeStr:  "0",
				expected: "status code must be between 100 and 511",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateStatusCode(tc.codeStr)

				assert.Error(t, err)
				assert.Equal(t, tc.expected, err.Error())
			})
		}
	})
}

func TestValidatePath(t *testing.T) {
	t.Run("happy path - valid paths", func(t *testing.T) {
		testCases := []struct {
			name string
			path string
		}{
			{
				name: "root path",
				path: "/",
			},
			{
				name: "simple path",
				path: "/api",
			},
			{
				name: "nested path",
				path: "/api/users",
			},
			{
				name: "path with trailing slash",
				path: "/api/users/",
			},
			{
				name: "path with query parameters",
				path: "/api/users?page=1",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidatePath(tc.path)

				assert.NoError(t, err)
			})
		}
	})

	t.Run("error path - path does not start with slash", func(t *testing.T) {
		testCases := []struct {
			name     string
			path     string
			expected string
		}{
			{
				name:     "path without leading slash",
				path:     "api/users",
				expected: "route path must start with '/'",
			},
			{
				name:     "empty string",
				path:     "",
				expected: "route path must start with '/'",
			},
			{
				name:     "path starting with letter",
				path:     "users",
				expected: "route path must start with '/'",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidatePath(tc.path)

				assert.Error(t, err)
				assert.Equal(t, tc.expected, err.Error())
			})
		}
	})
}

func TestValidateHost(t *testing.T) {
	t.Run("happy path - valid hosts", func(t *testing.T) {
		testCases := []struct {
			name string
			host string
		}{
			{
				name: "http URL",
				host: "http://localhost:8080",
			},
			{
				name: "https URL",
				host: "https://example.com",
			},
			{
				name: "URL with path",
				host: "http://localhost:8080/api",
			},
			{
				name: "URL with query parameters",
				host: "http://localhost:8080/api?key=value",
			},
			{
				name: "domain name",
				host: "example.com",
			},
			{
				name: "localhost",
				host: "localhost",
			},
			{
				name: "IP address",
				host: "127.0.0.1",
			},
			{
				name: "IP address with port in URL",
				host: "http://127.0.0.1:8080",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateHost(tc.host)

				assert.NoError(t, err)
			})
		}
	})

	t.Run("error path - invalid hosts", func(t *testing.T) {
		testCases := []struct {
			name     string
			host     string
			expected string
		}{
			{
				name:     "IP address with port without scheme",
				host:     "127.0.0.1:8080",
				expected: "invalid hostname",
			},
			{
				name:     "malformed URL - missing scheme",
				host:     "://invalid",
				expected: "invalid hostname",
			},
			{
				name:     "URL with invalid bracket",
				host:     "http://[invalid",
				expected: "invalid hostname",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateHost(tc.host)

				assert.Error(t, err)
				assert.Equal(t, tc.expected, err.Error())
			})
		}
	})
}

func TestValidateNonEmpty(t *testing.T) {
	t.Run("happy path - non-empty values", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     string
			fieldName string
		}{
			{
				name:      "simple string",
				value:     "test",
				fieldName: "field",
			},
			{
				name:      "single character",
				value:     "a",
				fieldName: "name",
			},
			{
				name:      "string with spaces",
				value:     "test value",
				fieldName: "description",
			},
			{
				name:      "numeric string",
				value:     "123",
				fieldName: "id",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateNonEmpty(tc.value, tc.fieldName)

				assert.NoError(t, err)
			})
		}
	})

	t.Run("error path - empty values", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     string
			fieldName string
			expected  string
		}{
			{
				name:      "empty string",
				value:     "",
				fieldName: "name",
				expected:  "name cannot be empty",
			},
			{
				name:      "empty string with different field name",
				value:     "",
				fieldName: "description",
				expected:  "description cannot be empty",
			},
			{
				name:      "empty string with custom field name",
				value:     "",
				fieldName: "host",
				expected:  "host cannot be empty",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateNonEmpty(tc.value, tc.fieldName)

				assert.Error(t, err)
				assert.Equal(t, tc.expected, err.Error())
			})
		}
	})
}
