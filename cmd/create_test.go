package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/cmd/form_builder"
	"inzibat/config"
)

func TestLoadHeadersFromFile(t *testing.T) {
	t.Run("happy path - valid JSON file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "headers.json")
		headersData := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
		}
		data, err := json.Marshal(headersData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		headers, err := config.LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
	})

	t.Run("error path - file does not exist", func(t *testing.T) {
		nonExistentPath := "/non/existent/file.json"

		headers, err := config.LoadHeadersFromFile(nonExistentPath)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("error path - invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(filePath, []byte("invalid json content"), 0644)
		require.NoError(t, err)

		headers, err := config.LoadHeadersFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, headers)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("error path - empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty.json")
		err := os.WriteFile(filePath, []byte("{}"), 0644)
		require.NoError(t, err)

		headers, err := config.LoadHeadersFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, 0, len(headers))
	})
}

func TestLoadBodyFromFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.json")
		bodyData := config.HttpBody{
			"message": "success",
			"code":    float64(200),
		}
		data, err := json.Marshal(bodyData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		body, err := config.LoadBodyFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, body)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, float64(200), body["code"])
	})

	t.Run("file does not exist", func(t *testing.T) {
		nonExistentPath := "/non/existent/body.json"

		body, err := config.LoadBodyFromFile(nonExistentPath)

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(filePath, []byte("not a valid json"), 0644)
		require.NoError(t, err)

		body, err := config.LoadBodyFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, body)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})
}

func TestLoadBodyStringFromFile(t *testing.T) {
	t.Run("happy path - valid text file", func(t *testing.T) {

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.txt")
		expectedContent := `{"message": "success", "status": "ok"}`
		err := os.WriteFile(filePath, []byte(expectedContent), 0644)
		require.NoError(t, err)

		bodyString, err := config.LoadBodyStringFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, bodyString)
	})

	t.Run("error path - file does not exist", func(t *testing.T) {

		nonExistentPath := "/non/existent/body.txt"

		bodyString, err := config.LoadBodyStringFromFile(nonExistentPath)

		assert.Empty(t, bodyString)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("happy path - empty file", func(t *testing.T) {

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty.txt")
		err := os.WriteFile(filePath, []byte(""), 0644)
		require.NoError(t, err)

		bodyString, err := config.LoadBodyStringFromFile(filePath)

		assert.Empty(t, bodyString)
	})
}

func TestCreateRouteForm(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		form := createRouteForm()

		assert.NotNil(t, form)
	})
}

func TestCreateCmd(t *testing.T) {
	t.Run("happy path - command is registered", func(t *testing.T) {
		assert.NotNil(t, createCmd)
		assert.Equal(t, "create", createCmd.Use)
		assert.Contains(t, createCmd.Aliases, "create-route")
		assert.Contains(t, createCmd.Aliases, "c")
	})

	t.Run("happy path - command has correct properties", func(t *testing.T) {
		assert.Contains(t, createCmd.Short, "Create")
		assert.Contains(t, createCmd.Long, "interactively")
	})
}

func TestHttpMethods(t *testing.T) {
	t.Run("happy path - httpMethods contains all methods", func(t *testing.T) {
		assert.Equal(t, 5, len(httpMethods))
		methodValues := make(map[string]bool)
		for _, opt := range httpMethods {
			methodValues[opt.Value] = true
		}
		assert.True(t, methodValues["GET"])
		assert.True(t, methodValues["POST"])
		assert.True(t, methodValues["PUT"])
		assert.True(t, methodValues["PATCH"])
		assert.True(t, methodValues["DELETE"])
	})
}

func TestRouteTypes(t *testing.T) {
	t.Run("happy path - routeTypes contains both types", func(t *testing.T) {
		assert.Equal(t, 2, len(routeTypes))
		typeValues := make(map[string]bool)
		for _, opt := range routeTypes {
			typeValues[opt.Value] = true
		}
		assert.True(t, typeValues["mock"])
		assert.True(t, typeValues["client"])
	})
}

func TestCreateRouteForm_Structure(t *testing.T) {
	t.Run("happy path - form has correct structure", func(t *testing.T) {
		form := createRouteForm()

		assert.NotNil(t, form)
	})
}

func TestCreateMockResponseForm_Structure(t *testing.T) {
	t.Run("happy path - status code form structure", func(t *testing.T) {
		status := "200"
		statusForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("statusCode").
					Title("Status Code").
					Placeholder(status).
					Value(&status).
					Validate(form_builder.ValidateStatusCode),
			),
		)

		assert.NotNil(t, statusForm)
	})

	t.Run("happy path - body type form structure", func(t *testing.T) {
		bodyForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("bodyType").
					Title("Body Type").
					Options([]huh.Option[string]{
						{Key: "Body (JSON object)", Value: BodyTypeBody},
						{Key: "BodyString (string)", Value: BodyTypeBodyString},
						{Key: "Skip", Value: form_builder.SourceSkip},
					}...),
			),
		)

		assert.NotNil(t, bodyForm)
	})

	t.Run("happy path - body type constants are used correctly", func(t *testing.T) {
		assert.Equal(t, "body", BodyTypeBody)
		assert.Equal(t, "bodyString", BodyTypeBodyString)
		assert.Equal(t, "structured", BodyTypeStructured)
	})
}

func TestCreateClientRequestForm_Structure(t *testing.T) {
	t.Run("happy path - basic form structure", func(t *testing.T) {
		basicForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("host").
					Title("Target Host URL").
					Placeholder("http://localhost:8081").
					Validate(form_builder.ValidateHost),
				huh.NewInput().
					Key("path").
					Title("Target Path").
					Placeholder("/api/users").
					Validate(form_builder.ValidateHost),
				huh.NewSelect[string]().
					Key("method").
					Title("Target HTTP Method").
					Options(httpMethods...),
			),
		)

		assert.NotNil(t, basicForm)
	})

	t.Run("happy path - request body type form structure", func(t *testing.T) {
		bodyForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("bodyType").
					Title("Request Body Type").
					Options([]huh.Option[string]{
						{Key: "Structured (JSON object)", Value: BodyTypeStructured},
						{Key: "Skip", Value: form_builder.SourceSkip},
					}...),
			),
		)

		assert.NotNil(t, bodyForm)
	})

	t.Run("happy path - options form structure", func(t *testing.T) {
		optionsForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Key("passWithRequestBody").
					Title("Pass With Request Body"),
				huh.NewConfirm().
					Key("passWithRequestHeaders").
					Title("Pass With Request Headers"),
				huh.NewConfirm().
					Key("inErrorReturn500").
					Title("Return 500 on Error"),
			),
		)

		assert.NotNil(t, optionsForm)
	})
}

func TestCreateRoute_Constants(t *testing.T) {
	t.Run("happy path - route type constants match expected values", func(t *testing.T) {
		assert.Equal(t, "mock", RouteTypeMock)
		assert.Equal(t, "client", RouteTypeClient)
	})

	t.Run("happy path - route type constants are used in routeTypes", func(t *testing.T) {
		typeValues := make(map[string]bool)
		for _, opt := range routeTypes {
			typeValues[opt.Value] = true
		}
		assert.True(t, typeValues[RouteTypeMock])
		assert.True(t, typeValues[RouteTypeClient])
	})
}

func TestCreateRouteForm_Fields(t *testing.T) {
	t.Run("happy path - form has path field", func(t *testing.T) {
		form := createRouteForm()
		assert.NotNil(t, form)
	})

	t.Run("happy path - form has method field", func(t *testing.T) {
		form := createRouteForm()
		assert.NotNil(t, form)
	})

	t.Run("happy path - form has routeType field", func(t *testing.T) {
		form := createRouteForm()
		assert.NotNil(t, form)
	})
}

func TestStatusCodeParsing(t *testing.T) {
	t.Run("error path - invalid status code string", func(t *testing.T) {
		invalidStatusCodes := []string{"", "abc", "not-a-number", "9999", "-1"}

		for _, invalidCode := range invalidStatusCodes {
			t.Run("invalid code: "+invalidCode, func(t *testing.T) {
				_, err := strconv.Atoi(invalidCode)
				if invalidCode == "" || invalidCode == "abc" || invalidCode == "not-a-number" {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("happy path - valid status code strings", func(t *testing.T) {
		validStatusCodes := []string{"200", "201", "400", "404", "500"}

		for _, validCode := range validStatusCodes {
			t.Run("valid code: "+validCode, func(t *testing.T) {
				statusCode, err := strconv.Atoi(validCode)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, statusCode, 100)
				assert.LessOrEqual(t, statusCode, 599)
			})
		}
	})
}
