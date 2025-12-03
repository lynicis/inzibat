package config

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCfg_GetServerAddr(t *testing.T) {
	t.Run("happy path - returns formatted server address", func(t *testing.T) {
		cfg := &Cfg{
			ServerPort: 8080,
		}

		addr := cfg.GetServerAddr()

		assert.Equal(t, ":8080", addr)
	})

	t.Run("happy path - handles different ports", func(t *testing.T) {
		testCases := []struct {
			port     int
			expected string
		}{
			{3000, ":3000"},
			{443, ":443"},
			{80, ":80"},
			{9999, ":9999"},
		}

		for _, tc := range testCases {
			cfg := &Cfg{ServerPort: tc.port}
			assert.Equal(t, tc.expected, cfg.GetServerAddr())
		}
	})
}

func TestCfg_ConvertRoutesTuiTable(t *testing.T) {
	t.Run("happy path - converts routes to table rows with correct types", func(t *testing.T) {

		cfg := &Cfg{
			Routes: []Route{
				{
					Method: "GET",
					Path:   "/mock",
					FakeResponse: &FakeResponse{
						StatusCode: 200,
					},
				},
				{
					Method: "POST",
					Path:   "/proxy",
					RequestTo: &RequestTo{
						Host: "http://example.com",
						Path: "/proxy",
					},
				},
				{
					Method: "DELETE",
					Path:   "/unknown",
				},
			},
		}

		rows := cfg.ConvertRoutesTuiTable()

		assert.Len(t, rows, 3)

		assert.Equal(t, []string{"GET", "/mock", "MOCK"}, rows[0])

		assert.Equal(t, []string{"POST", "/proxy", "PROXY"}, rows[1])

		assert.Equal(t, []string{"DELETE", "/unknown", "UNKNOWN"}, rows[2])
	})
}

func TestRequestTo_GetParsedUrl(t *testing.T) {
	t.Run("happy path - parses valid URL", func(t *testing.T) {
		requestTo := &RequestTo{
			Host: "http://localhost",
			Path: "/api/users",
		}

		parsedURL, err := requestTo.GetParsedUrl()

		assert.NoError(t, err)
		assert.NotNil(t, parsedURL)
		assert.Equal(t, "http://localhost/api/users", parsedURL.String())
	})

	t.Run("happy path - parses URL with port", func(t *testing.T) {
		requestTo := &RequestTo{
			Host: "http://localhost:8080",
			Path: "/test",
		}

		parsedURL, err := requestTo.GetParsedUrl()

		assert.NoError(t, err)
		assert.NotNil(t, parsedURL)
		assert.Equal(t, "http://localhost:8080/test", parsedURL.String())
	})

	t.Run("happy path - parses HTTPS URL", func(t *testing.T) {
		requestTo := &RequestTo{
			Host: "https://api.example.com",
			Path: "/v1/data",
		}

		parsedURL, err := requestTo.GetParsedUrl()

		assert.NoError(t, err)
		assert.NotNil(t, parsedURL)
		assert.Equal(t, "https://api.example.com/v1/data", parsedURL.String())
	})

	t.Run("error path - invalid URL", func(t *testing.T) {
		requestTo := &RequestTo{
			Host: "://invalid",
			Path: "/test",
		}

		parsedURL, err := requestTo.GetParsedUrl()

		assert.Error(t, err)
		assert.Nil(t, parsedURL)
		assert.Contains(t, err.Error(), "failed to parse url")
	})

	t.Run("happy path - empty path", func(t *testing.T) {
		requestTo := &RequestTo{
			Host: "http://localhost",
			Path: "/",
		}

		parsedURL, err := requestTo.GetParsedUrl()

		assert.NoError(t, err)
		assert.NotNil(t, parsedURL)
		expectedURL, _ := url.Parse("http://localhost/")
		assert.Equal(t, expectedURL.String(), parsedURL.String())
	})
}
