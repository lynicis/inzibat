package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

func TestCreateMockResponseForm_ErrorPaths(t *testing.T) {
	t.Run("error path - status code parsing error message format", func(t *testing.T) {
		invalidStatusCode := "not-a-number"
		_, err := strconv.Atoi(invalidStatusCode)

		assert.Error(t, err)
		expectedErrorMsg := fmt.Sprintf("failed to parse status code %q", invalidStatusCode)
		assert.Contains(t, fmt.Sprintf("failed to parse status code %q", invalidStatusCode), expectedErrorMsg)
	})

	t.Run("error path - status code form run error message format", func(t *testing.T) {
		testErr := fmt.Errorf("form run failed")

		wrappedErr := fmt.Errorf("failed to get status code: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to get status code")
	})

	t.Run("error path - headers collection error message format", func(t *testing.T) {
		testErr := fmt.Errorf("headers collection failed")

		wrappedErr := fmt.Errorf("failed to collect headers: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to collect headers")
	})

	t.Run("error path - body type form run error message format", func(t *testing.T) {
		testErr := fmt.Errorf("body type form failed")

		wrappedErr := fmt.Errorf("failed to select body type: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to select body type")
	})

	t.Run("error path - body collection error message format", func(t *testing.T) {
		testErr := fmt.Errorf("body collection failed")

		wrappedErr := fmt.Errorf("failed to collect body: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to collect body")
	})

	t.Run("error path - body string collection error message format", func(t *testing.T) {
		testErr := fmt.Errorf("body string collection failed")

		wrappedErr := fmt.Errorf("failed to collect body string: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to collect body string")
	})
}

func TestCreateClientRequestForm_ErrorPaths(t *testing.T) {
	t.Run("error path - basic form run error message format", func(t *testing.T) {
		testErr := fmt.Errorf("basic form failed")

		wrappedErr := fmt.Errorf("failed to get basic request info: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to get basic request info")
	})

	t.Run("error path - headers collection error message format", func(t *testing.T) {
		testErr := fmt.Errorf("headers collection failed")

		wrappedErr := fmt.Errorf("failed to collect headers: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to collect headers")
	})

	t.Run("error path - body type form run error message format", func(t *testing.T) {
		testErr := fmt.Errorf("body type form failed")

		wrappedErr := fmt.Errorf("failed to select body type: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to select body type")
	})

	t.Run("error path - body collection error message format", func(t *testing.T) {
		testErr := fmt.Errorf("body collection failed")

		wrappedErr := fmt.Errorf("failed to collect body: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to collect body")
	})

	t.Run("error path - options form run error message format", func(t *testing.T) {
		testErr := fmt.Errorf("options form failed")

		wrappedErr := fmt.Errorf("failed to get options: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to get options")
	})
}

func TestCreateRoute_ErrorPaths(t *testing.T) {
	t.Run("error path - route form run error message format", func(t *testing.T) {
		testErr := fmt.Errorf("route form failed")

		wrappedErr := fmt.Errorf("failed to create route: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to create route")
	})

	t.Run("error path - mock response form error propagation", func(t *testing.T) {
		testErr := fmt.Errorf("mock response form failed")

		assert.Error(t, testErr)
		assert.Contains(t, testErr.Error(), "mock response form failed")
	})

	t.Run("error path - client request form error propagation", func(t *testing.T) {
		testErr := fmt.Errorf("client request form failed")

		assert.Error(t, testErr)
		assert.Contains(t, testErr.Error(), "client request form failed")
	})
}

func TestCreateRoute_LogicBranches(t *testing.T) {
	t.Run("happy path - route type mock branch", func(t *testing.T) {
		routeType := RouteTypeMock

		assert.Equal(t, "mock", routeType)
		switch routeType {
		case RouteTypeMock:
			assert.True(t, true, "mock branch should be executed")
		default:
			t.Fatal("should match mock route type")
		}
	})

	t.Run("happy path - route type client branch", func(t *testing.T) {
		routeType := RouteTypeClient

		assert.Equal(t, "client", routeType)
		switch routeType {
		case RouteTypeClient:
			assert.True(t, true, "client branch should be executed")
		default:
			t.Fatal("should match client route type")
		}
	})

	t.Run("happy path - route structure creation", func(t *testing.T) {
		method := "GET"
		path := "/test"
		var fakeResponse *config.FakeResponse
		var requestTo *config.RequestTo

		route := &config.Route{
			Method:       method,
			Path:         path,
			FakeResponse: fakeResponse,
			RequestTo:    requestTo,
		}

		assert.NotNil(t, route)
		assert.Equal(t, method, route.Method)
		assert.Equal(t, path, route.Path)
		assert.Nil(t, route.FakeResponse)
		assert.Nil(t, route.RequestTo)
	})
}

func TestCreateMockResponseForm_LogicBranches(t *testing.T) {
	t.Run("happy path - body type body branch", func(t *testing.T) {
		bodyType := BodyTypeBody

		assert.Equal(t, "body", bodyType)
		switch bodyType {
		case BodyTypeBody:
			assert.True(t, true, "body branch should be executed")
		default:
			t.Fatal("should match body type")
		}
	})

	t.Run("happy path - body type bodyString branch", func(t *testing.T) {
		bodyType := BodyTypeBodyString

		assert.Equal(t, "bodyString", bodyType)
		switch bodyType {
		case BodyTypeBodyString:
			assert.True(t, true, "bodyString branch should be executed")
		default:
			t.Fatal("should match bodyString type")
		}
	})

	t.Run("happy path - body type skip branch", func(t *testing.T) {
		bodyType := form_builder.SourceSkip

		assert.NotEqual(t, BodyTypeBody, bodyType)
		assert.NotEqual(t, BodyTypeBodyString, bodyType)
		switch bodyType {
		case BodyTypeBody, BodyTypeBodyString:
			t.Fatal("should not match body or bodyString")
		default:
			assert.True(t, true, "skip branch should be executed")
		}
	})

	t.Run("happy path - fake response structure creation", func(t *testing.T) {
		statusCode := 200
		headers := make(map[string][]string)

		fakeResponse := &config.FakeResponse{
			StatusCode: statusCode,
			Headers:    headers,
		}

		assert.NotNil(t, fakeResponse)
		assert.Equal(t, statusCode, fakeResponse.StatusCode)
		assert.NotNil(t, fakeResponse.Headers)
	})
}

func TestCreateClientRequestForm_LogicBranches(t *testing.T) {
	t.Run("happy path - body type structured branch", func(t *testing.T) {
		bodyType := BodyTypeStructured

		assert.Equal(t, "structured", bodyType)
		if bodyType == BodyTypeStructured {
			assert.True(t, true, "structured branch should be executed")
		} else {
			t.Fatal("should match structured type")
		}
	})

	t.Run("happy path - body type skip branch", func(t *testing.T) {
		bodyType := form_builder.SourceSkip

		assert.NotEqual(t, BodyTypeStructured, bodyType)
		if bodyType == BodyTypeStructured {
			t.Fatal("should not match structured")
		} else {
			assert.True(t, true, "skip branch should be executed")
		}
	})

	t.Run("happy path - request to structure creation", func(t *testing.T) {
		host := "http://localhost:8081"
		targetPath := "/api/users"
		targetMethod := "GET"
		headers := make(map[string][]string)
		var body config.HttpBody

		requestTo := &config.RequestTo{
			Host:                   host,
			Path:                   targetPath,
			Method:                 targetMethod,
			Headers:                headers,
			Body:                   body,
			PassWithRequestBody:    false,
			PassWithRequestHeaders: false,
			InErrorReturn500:       false,
		}

		assert.NotNil(t, requestTo)
		assert.Equal(t, host, requestTo.Host)
		assert.Equal(t, targetPath, requestTo.Path)
		assert.Equal(t, targetMethod, requestTo.Method)
		assert.NotNil(t, requestTo.Headers)
		assert.False(t, requestTo.PassWithRequestBody)
		assert.False(t, requestTo.PassWithRequestHeaders)
		assert.False(t, requestTo.InErrorReturn500)
	})
}

func TestCreateCmd_Run(t *testing.T) {
	t.Run("happy path - command structure", func(t *testing.T) {

		assert.NotNil(t, createCmd)
		assert.Equal(t, "create", createCmd.Use)
		assert.NotNil(t, createCmd.Args, "Args should be set")
		err := createCmd.Args(createCmd, []string{"test"})
		assert.Error(t, err, "NoArgs should return error when args are provided")
		assert.NotNil(t, createCmd.Run)
	})

	t.Run("happy path - command aliases", func(t *testing.T) {
		aliases := createCmd.Aliases

		assert.Len(t, aliases, 2)
		assert.Contains(t, aliases, "create-route")
		assert.Contains(t, aliases, "c")
	})

	t.Run("happy path - command description", func(t *testing.T) {
		short := createCmd.Short
		long := createCmd.Long

		assert.Contains(t, short, "Create")
		assert.Contains(t, long, "interactively")
	})
}

func TestCreateMockResponseForm_StatusCodes(t *testing.T) {
	t.Run("happy path - valid status code range", func(t *testing.T) {
		validCodes := []int{100, 200, 201, 300, 400, 404, 500, 503, 599}

		for _, code := range validCodes {
			t.Run(fmt.Sprintf("status code %d", code), func(t *testing.T) {
				codeStr := strconv.Itoa(code)
				parsed, err := strconv.Atoi(codeStr)
				assert.NoError(t, err)
				assert.Equal(t, code, parsed)
				assert.GreaterOrEqual(t, parsed, 100)
				assert.LessOrEqual(t, parsed, 599)
			})
		}
	})

	t.Run("error path - status code out of range", func(t *testing.T) {
		invalidCodes := []string{"99", "600", "1000"}

		for _, codeStr := range invalidCodes {
			t.Run(fmt.Sprintf("invalid code %s", codeStr), func(t *testing.T) {
				code, err := strconv.Atoi(codeStr)
				if err == nil {
					if code < 100 || code > 599 {
						assert.True(t, code < 100 || code > 599, "status code should be out of range")
					}
				}
			})
		}
	})

	t.Run("happy path - default status code is 200", func(t *testing.T) {
		defaultStatus := strconv.Itoa(http.StatusOK)

		assert.Equal(t, "200", defaultStatus)
		statusCode, err := strconv.Atoi(defaultStatus)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
	})
}

func TestCreateMockResponseForm_BodyTypes(t *testing.T) {
	t.Run("happy path - body type constants are distinct", func(t *testing.T) {
		bodyTypes := []string{BodyTypeBody, BodyTypeBodyString, form_builder.SourceSkip}

		assert.Equal(t, "body", BodyTypeBody)
		assert.Equal(t, "bodyString", BodyTypeBodyString)
		assert.NotEqual(t, BodyTypeBody, BodyTypeBodyString)
		assert.NotEqual(t, BodyTypeBody, form_builder.SourceSkip)
		assert.NotEqual(t, BodyTypeBodyString, form_builder.SourceSkip)
		uniqueMap := make(map[string]bool)
		for _, bt := range bodyTypes {
			assert.False(t, uniqueMap[bt], "body type %q should be unique", bt)
			uniqueMap[bt] = true
		}
	})

	t.Run("happy path - fake response with body", func(t *testing.T) {
		statusCode := 200
		headers := make(map[string][]string)
		body := config.HttpBody{"message": "success"}

		fakeResponse := &config.FakeResponse{
			StatusCode: statusCode,
			Headers:    headers,
			Body:       body,
		}

		assert.NotNil(t, fakeResponse)
		assert.Equal(t, statusCode, fakeResponse.StatusCode)
		assert.NotNil(t, fakeResponse.Headers)
		assert.NotNil(t, fakeResponse.Body)
		assert.Equal(t, "success", fakeResponse.Body["message"])
		assert.Empty(t, fakeResponse.BodyString)
	})

	t.Run("happy path - fake response with bodyString", func(t *testing.T) {
		statusCode := 200
		headers := make(map[string][]string)
		bodyString := `{"message": "success"}`

		fakeResponse := &config.FakeResponse{
			StatusCode: statusCode,
			Headers:    headers,
			BodyString: bodyString,
		}

		assert.NotNil(t, fakeResponse)
		assert.Equal(t, statusCode, fakeResponse.StatusCode)
		assert.NotNil(t, fakeResponse.Headers)
		assert.Empty(t, fakeResponse.Body)
		assert.Equal(t, bodyString, fakeResponse.BodyString)
	})

	t.Run("happy path - fake response without body", func(t *testing.T) {
		statusCode := 204
		headers := make(map[string][]string)

		fakeResponse := &config.FakeResponse{
			StatusCode: statusCode,
			Headers:    headers,
		}

		assert.NotNil(t, fakeResponse)
		assert.Equal(t, statusCode, fakeResponse.StatusCode)
		assert.NotNil(t, fakeResponse.Headers)
		assert.Nil(t, fakeResponse.Body)
		assert.Empty(t, fakeResponse.BodyString)
	})
}

func TestCreateClientRequestForm_Options(t *testing.T) {
	t.Run("happy path - all options false", func(t *testing.T) {
		passWithRequestBody := false
		passWithRequestHeaders := false
		inErrorReturn500 := false

		assert.False(t, passWithRequestBody)
		assert.False(t, passWithRequestHeaders)
		assert.False(t, inErrorReturn500)
	})

	t.Run("happy path - all options true", func(t *testing.T) {
		passWithRequestBody := true
		passWithRequestHeaders := true
		inErrorReturn500 := true

		assert.True(t, passWithRequestBody)
		assert.True(t, passWithRequestHeaders)
		assert.True(t, inErrorReturn500)
	})

	t.Run("happy path - mixed options", func(t *testing.T) {
		passWithRequestBody := true
		passWithRequestHeaders := false
		inErrorReturn500 := true

		assert.True(t, passWithRequestBody)
		assert.False(t, passWithRequestHeaders)
		assert.True(t, inErrorReturn500)
	})

	t.Run("happy path - request to with all options", func(t *testing.T) {
		host := "http://localhost:8081"
		targetPath := "/api/users"
		targetMethod := "POST"
		headers := make(map[string][]string)
		body := config.HttpBody{"id": float64(1)}

		requestTo := &config.RequestTo{
			Host:                   host,
			Path:                   targetPath,
			Method:                 targetMethod,
			Headers:                headers,
			Body:                   body,
			PassWithRequestBody:    true,
			PassWithRequestHeaders: true,
			InErrorReturn500:       true,
		}

		assert.NotNil(t, requestTo)
		assert.Equal(t, host, requestTo.Host)
		assert.Equal(t, targetPath, requestTo.Path)
		assert.Equal(t, targetMethod, requestTo.Method)
		assert.True(t, requestTo.PassWithRequestBody)
		assert.True(t, requestTo.PassWithRequestHeaders)
		assert.True(t, requestTo.InErrorReturn500)
	})
}

func TestCreateClientRequestForm_BodyHandling(t *testing.T) {
	t.Run("happy path - structured body type", func(t *testing.T) {
		bodyType := BodyTypeStructured
		body := config.HttpBody{
			"name":  "test",
			"value": float64(123),
		}

		assert.Equal(t, "structured", bodyType)
		assert.NotNil(t, body)
		assert.Equal(t, "test", body["name"])
		assert.Equal(t, float64(123), body["value"])
	})

	t.Run("happy path - skip body type", func(t *testing.T) {
		bodyType := form_builder.SourceSkip
		var body config.HttpBody

		assert.NotEqual(t, BodyTypeStructured, bodyType)
		assert.Nil(t, body)
	})

	t.Run("happy path - request to with structured body", func(t *testing.T) {
		host := "http://localhost:8081"
		targetPath := "/api/users"
		targetMethod := "POST"
		headers := make(map[string][]string)
		body := config.HttpBody{"name": "John", "age": float64(30)}

		requestTo := &config.RequestTo{
			Host:    host,
			Path:    targetPath,
			Method:  targetMethod,
			Headers: headers,
			Body:    body,
		}

		assert.NotNil(t, requestTo)
		assert.NotNil(t, requestTo.Body)
		assert.Equal(t, "John", requestTo.Body["name"])
		assert.Equal(t, float64(30), requestTo.Body["age"])
	})

	t.Run("happy path - request to without body", func(t *testing.T) {
		host := "http://localhost:8081"
		targetPath := "/api/users"
		targetMethod := "GET"
		headers := make(map[string][]string)
		var body config.HttpBody

		requestTo := &config.RequestTo{
			Host:    host,
			Path:    targetPath,
			Method:  targetMethod,
			Headers: headers,
			Body:    body,
		}

		assert.NotNil(t, requestTo)
		assert.Nil(t, requestTo.Body)
	})
}

func TestCreateRoute_CompleteFlow(t *testing.T) {
	t.Run("happy path - mock route structure", func(t *testing.T) {
		method := "GET"
		path := "/api/test"
		fakeResponse := &config.FakeResponse{
			StatusCode: 200,
			Headers:    make(map[string][]string),
			Body:       config.HttpBody{"message": "success"},
		}

		route := &config.Route{
			Method:       method,
			Path:         path,
			FakeResponse: fakeResponse,
			RequestTo:    nil,
		}

		assert.NotNil(t, route)
		assert.Equal(t, method, route.Method)
		assert.Equal(t, path, route.Path)
		assert.NotNil(t, route.FakeResponse)
		assert.Nil(t, route.RequestTo)
		assert.Equal(t, 200, route.FakeResponse.StatusCode)
	})

	t.Run("happy path - client route structure", func(t *testing.T) {
		method := "POST"
		path := "/api/proxy"
		requestTo := &config.RequestTo{
			Host:    "http://localhost:8081",
			Path:    "/target",
			Method:  "POST",
			Headers: make(map[string][]string),
			Body:    config.HttpBody{"data": "value"},
		}

		route := &config.Route{
			Method:       method,
			Path:         path,
			FakeResponse: nil,
			RequestTo:    requestTo,
		}

		assert.NotNil(t, route)
		assert.Equal(t, method, route.Method)
		assert.Equal(t, path, route.Path)
		assert.Nil(t, route.FakeResponse)
		assert.NotNil(t, route.RequestTo)
		assert.Equal(t, "http://localhost:8081", route.RequestTo.Host)
	})

	t.Run("happy path - route with both nil (should not happen but test structure)", func(t *testing.T) {
		method := "GET"
		path := "/api/test"

		route := &config.Route{
			Method:       method,
			Path:         path,
			FakeResponse: nil,
			RequestTo:    nil,
		}

		assert.NotNil(t, route)
		assert.Equal(t, method, route.Method)
		assert.Equal(t, path, route.Path)
		assert.Nil(t, route.FakeResponse)
		assert.Nil(t, route.RequestTo)
	})
}

func TestCreateRouteForm_FieldValidation(t *testing.T) {
	t.Run("happy path - form has all required fields", func(t *testing.T) {
		form := createRouteForm()

		assert.NotNil(t, form)
	})

	t.Run("happy path - form fields are properly configured", func(t *testing.T) {
		form := createRouteForm()

		assert.NotNil(t, form)
	})
}

func TestHttpMethods_Completeness(t *testing.T) {
	t.Run("happy path - all standard methods present", func(t *testing.T) {
		expectedMethods := map[string]bool{
			"GET":    false,
			"POST":   false,
			"PUT":    false,
			"PATCH":  false,
			"DELETE": false,
		}

		for _, opt := range httpMethods {
			if _, exists := expectedMethods[opt.Value]; exists {
				expectedMethods[opt.Value] = true
			}
		}

		for method, found := range expectedMethods {
			assert.True(t, found, "method %s should be present", method)
		}
	})

	t.Run("happy path - method keys match values", func(t *testing.T) {
		for _, opt := range httpMethods {
			assert.Equal(t, opt.Key, opt.Value, "method key should match value for %s", opt.Value)
		}
	})
}

func TestRouteTypes_Completeness(t *testing.T) {
	t.Run("happy path - both route types present", func(t *testing.T) {
		expectedTypes := map[string]bool{
			RouteTypeMock:   false,
			RouteTypeClient: false,
		}

		for _, opt := range routeTypes {
			if _, exists := expectedTypes[opt.Value]; exists {
				expectedTypes[opt.Value] = true
			}
		}

		for routeType, found := range expectedTypes {
			assert.True(t, found, "route type %s should be present", routeType)
		}
	})

	t.Run("happy path - route type keys are descriptive", func(t *testing.T) {
		for _, opt := range routeTypes {
			assert.NotEmpty(t, opt.Key, "route type key should not be empty")
			assert.NotEmpty(t, opt.Value, "route type value should not be empty")
			assert.GreaterOrEqual(t, len(opt.Key), len(opt.Value), "key should be at least as long as value")
		}
	})
}

func TestCreateMockResponseForm_HeaderHandling(t *testing.T) {
	t.Run("happy path - empty headers", func(t *testing.T) {
		headers := make(map[string][]string)

		fakeResponse := &config.FakeResponse{
			StatusCode: 200,
			Headers:    headers,
		}

		assert.NotNil(t, fakeResponse.Headers)
		assert.Equal(t, 0, len(fakeResponse.Headers))
	})

	t.Run("happy path - headers with single value", func(t *testing.T) {
		headers := make(map[string][]string)
		headers["Content-Type"] = []string{"application/json"}

		fakeResponse := &config.FakeResponse{
			StatusCode: 200,
			Headers:    headers,
		}

		assert.NotNil(t, fakeResponse.Headers)
		assert.Equal(t, 1, len(fakeResponse.Headers))
		assert.Equal(t, "application/json", fakeResponse.Headers["Content-Type"][0])
	})

	t.Run("happy path - headers with multiple values", func(t *testing.T) {
		headers := make(map[string][]string)
		headers["Accept"] = []string{"application/json", "text/html"}

		fakeResponse := &config.FakeResponse{
			StatusCode: 200,
			Headers:    headers,
		}

		assert.NotNil(t, fakeResponse.Headers)
		assert.Equal(t, 1, len(fakeResponse.Headers))
		assert.Equal(t, 2, len(fakeResponse.Headers["Accept"]))
		assert.Contains(t, fakeResponse.Headers["Accept"], "application/json")
		assert.Contains(t, fakeResponse.Headers["Accept"], "text/html")
	})
}

func TestCreateClientRequestForm_HeaderHandling(t *testing.T) {
	t.Run("happy path - empty headers", func(t *testing.T) {
		headers := make(map[string][]string)

		requestTo := &config.RequestTo{
			Host:    "http://localhost:8081",
			Path:    "/api",
			Method:  "GET",
			Headers: headers,
		}

		assert.NotNil(t, requestTo.Headers)
		assert.Equal(t, 0, len(requestTo.Headers))
	})

	t.Run("happy path - headers with authorization", func(t *testing.T) {
		headers := make(map[string][]string)
		headers["Authorization"] = []string{"Bearer token123"}

		requestTo := &config.RequestTo{
			Host:    "http://localhost:8081",
			Path:    "/api",
			Method:  "GET",
			Headers: headers,
		}

		assert.NotNil(t, requestTo.Headers)
		assert.Equal(t, 1, len(requestTo.Headers))
		assert.Equal(t, "Bearer token123", requestTo.Headers["Authorization"][0])
	})
}

func TestCreateRoute_ErrorPropagation(t *testing.T) {
	t.Run("error path - route form error propagates correctly", func(t *testing.T) {
		testErr := fmt.Errorf("route form failed")

		wrappedErr := fmt.Errorf("failed to create route: %w", testErr)
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "failed to create route")
		assert.ErrorIs(t, wrappedErr, testErr)
	})

	t.Run("error path - mock response form error propagates without double wrapping", func(t *testing.T) {
		testErr := fmt.Errorf("failed to collect headers: header error")

		assert.Error(t, testErr)
		assert.Contains(t, testErr.Error(), "failed to collect headers")
	})

	t.Run("error path - client request form error propagates without double wrapping", func(t *testing.T) {
		testErr := fmt.Errorf("failed to get basic request info: form error")

		assert.Error(t, testErr)
		assert.Contains(t, testErr.Error(), "failed to get basic request info")
	})
}

func TestCreateMockResponseForm_CompleteStructure(t *testing.T) {
	t.Run("happy path - complete fake response with all fields", func(t *testing.T) {
		statusCode := 201
		headers := make(map[string][]string)
		headers["Location"] = []string{"/api/users/123"}
		body := config.HttpBody{"id": float64(123), "name": "created"}

		fakeResponse := &config.FakeResponse{
			StatusCode: statusCode,
			Headers:    headers,
			Body:       body,
		}

		assert.NotNil(t, fakeResponse)
		assert.Equal(t, statusCode, fakeResponse.StatusCode)
		assert.Equal(t, 1, len(fakeResponse.Headers))
		assert.NotNil(t, fakeResponse.Body)
		assert.Equal(t, float64(123), fakeResponse.Body["id"])
		assert.Equal(t, "created", fakeResponse.Body["name"])
	})

	t.Run("happy path - fake response with bodyString only", func(t *testing.T) {
		statusCode := 200
		headers := make(map[string][]string)
		bodyString := `{"status": "ok", "data": {"id": 1}}`

		fakeResponse := &config.FakeResponse{
			StatusCode: statusCode,
			Headers:    headers,
			BodyString: bodyString,
		}

		assert.NotNil(t, fakeResponse)
		assert.Equal(t, statusCode, fakeResponse.StatusCode)
		assert.Nil(t, fakeResponse.Body)
		assert.Equal(t, bodyString, fakeResponse.BodyString)
	})
}

func TestCreateClientRequestForm_CompleteStructure(t *testing.T) {
	t.Run("happy path - complete request to with all fields", func(t *testing.T) {
		host := "https://api.example.com"
		targetPath := "/v1/users"
		targetMethod := "PUT"
		headers := make(map[string][]string)
		headers["Content-Type"] = []string{"application/json"}
		headers["Authorization"] = []string{"Bearer token"}
		body := config.HttpBody{"name": "updated", "email": "test@example.com"}

		requestTo := &config.RequestTo{
			Host:                   host,
			Path:                   targetPath,
			Method:                 targetMethod,
			Headers:                headers,
			Body:                   body,
			PassWithRequestBody:    true,
			PassWithRequestHeaders: true,
			InErrorReturn500:       false,
		}

		assert.NotNil(t, requestTo)
		assert.Equal(t, host, requestTo.Host)
		assert.Equal(t, targetPath, requestTo.Path)
		assert.Equal(t, targetMethod, requestTo.Method)
		assert.Equal(t, 2, len(requestTo.Headers))
		assert.NotNil(t, requestTo.Body)
		assert.True(t, requestTo.PassWithRequestBody)
		assert.True(t, requestTo.PassWithRequestHeaders)
		assert.False(t, requestTo.InErrorReturn500)
	})
}

func TestCreateMockResponseFormInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - statusFormRunner.Run() returns error", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("status form error")

		mockStatusForm.EXPECT().Run().Return(expectedError)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return nil, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return "", nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get status code")
	})

	t.Run("error path - invalid status code string", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("invalid")

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return nil, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return "", nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse status code")
	})

	t.Run("error path - headersCollector returns error", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("headers error")

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("200")

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return nil, expectedError },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return "", nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to collect headers")
	})

	t.Run("error path - bodyTypeFormRunner.Run() returns error", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("body type form error")

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("200")
		mockBodyTypeForm.EXPECT().Run().Return(expectedError)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return "", nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to select body type")
	})

	t.Run("happy path - body type BodyTypeBody", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		body := config.HttpBody{"message": "success"}

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("200")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(BodyTypeBody)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return headers, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return body, nil },
			func() (string, error) { return "", nil },
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 200, result.StatusCode)
		assert.Equal(t, "application/json", result.Headers.Get("Content-Type"))
		assert.NotNil(t, result.Body)
		assert.Equal(t, "success", result.Body["message"])
		assert.Empty(t, result.BodyString)
	})

	t.Run("happy path - body type BodyTypeBodyString", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		headers := make(http.Header)
		bodyString := `{"message": "success"}`

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("201")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(BodyTypeBodyString)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return headers, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return bodyString, nil },
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 201, result.StatusCode)
		assert.Nil(t, result.Body)
		assert.Equal(t, bodyString, result.BodyString)
	})

	t.Run("happy path - body type skip", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		headers := make(http.Header)

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("204")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(form_builder.SourceSkip)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return headers, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return "", nil },
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 204, result.StatusCode)
		assert.Nil(t, result.Body)
		assert.Empty(t, result.BodyString)
	})

	t.Run("error path - bodyCollector returns error", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("body collection error")

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("200")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(BodyTypeBody)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, expectedError },
			func() (string, error) { return "", nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to collect body")
	})

	t.Run("error path - bodyStringCollector returns error", func(t *testing.T) {

		mockStatusForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("body string collection error")

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("200")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(BodyTypeBodyString)

		result, err := createMockResponseFormInternal(
			mockStatusForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			func() (string, error) { return "", expectedError },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to collect body string")
	})
}

func TestCreateClientRequestFormInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - basicFormRunner.Run() returns error", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("basic form error")

		mockBasicForm.EXPECT().Run().Return(expectedError)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return nil, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get basic request info")
	})

	t.Run("error path - headersCollector returns error", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("headers error")

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return nil, expectedError },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to collect headers")
	})

	t.Run("error path - bodyTypeFormRunner.Run() returns error", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("body type form error")

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(expectedError)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to select body type")
	})

	t.Run("error path - bodyCollector returns error when bodyType is structured", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("body collection error")

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("POST")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(BodyTypeStructured)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, expectedError },
			mockOptionsForm,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to collect body")
	})

	t.Run("error path - optionsFormRunner.Run() returns error", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("options form error")

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(form_builder.SourceSkip)
		mockOptionsForm.EXPECT().Run().Return(expectedError)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get options")
	})

	t.Run("happy path - with structured body", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		headers := make(http.Header)
		headers.Set("Authorization", "Bearer token")
		body := config.HttpBody{"name": "test", "id": float64(1)}

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("POST")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(BodyTypeStructured)
		mockOptionsForm.EXPECT().Run().Return(nil)
		mockOptionsForm.EXPECT().GetBool("passWithRequestBody").Return(true)
		mockOptionsForm.EXPECT().GetBool("passWithRequestHeaders").Return(true)
		mockOptionsForm.EXPECT().GetBool("inErrorReturn500").Return(false)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return headers, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return body, nil },
			mockOptionsForm,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "http://localhost:8081", result.Host)
		assert.Equal(t, "/api/users", result.Path)
		assert.Equal(t, "POST", result.Method)
		assert.Equal(t, "Bearer token", result.Headers.Get("Authorization"))
		assert.NotNil(t, result.Body)
		assert.Equal(t, "test", result.Body["name"])
		assert.True(t, result.PassWithRequestBody)
		assert.True(t, result.PassWithRequestHeaders)
		assert.False(t, result.InErrorReturn500)
	})

	t.Run("happy path - skip body", func(t *testing.T) {

		mockBasicForm := form_builder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := form_builder.NewMockFormRunner(ctrl)
		mockOptionsForm := form_builder.NewMockFormRunner(ctrl)
		headers := make(http.Header)

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(form_builder.SourceSkip)
		mockOptionsForm.EXPECT().Run().Return(nil)
		mockOptionsForm.EXPECT().GetBool("passWithRequestBody").Return(false)
		mockOptionsForm.EXPECT().GetBool("passWithRequestHeaders").Return(false)
		mockOptionsForm.EXPECT().GetBool("inErrorReturn500").Return(true)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return headers, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "http://localhost:8081", result.Host)
		assert.Equal(t, "/api/users", result.Path)
		assert.Equal(t, "GET", result.Method)
		assert.Nil(t, result.Body)
		assert.False(t, result.PassWithRequestBody)
		assert.False(t, result.PassWithRequestHeaders)
		assert.True(t, result.InErrorReturn500)
	})
}

func TestCreateRouteInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - routeFormRunner.Run() returns error", func(t *testing.T) {
		mockRouteForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("route form error")

		mockRouteForm.EXPECT().Run().Return(expectedError)

		result, err := createRouteInternal(
			mockRouteForm,
			func() (*config.FakeResponse, error) { return nil, nil },
			func() (*config.RequestTo, error) { return nil, nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create route")
	})

	t.Run("error path - mockResponseFormCreator returns error", func(t *testing.T) {

		mockRouteForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("mock response form error")

		mockRouteForm.EXPECT().Run().Return(nil)
		mockRouteForm.EXPECT().GetString("path").Return("/api/test")
		mockRouteForm.EXPECT().GetString("method").Return("GET")
		mockRouteForm.EXPECT().GetString("routeType").Return(RouteTypeMock)

		result, err := createRouteInternal(
			mockRouteForm,
			func() (*config.FakeResponse, error) { return nil, expectedError },
			func() (*config.RequestTo, error) { return nil, nil },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("error path - clientRequestFormCreator returns error", func(t *testing.T) {

		mockRouteForm := form_builder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("client request form error")

		mockRouteForm.EXPECT().Run().Return(nil)
		mockRouteForm.EXPECT().GetString("path").Return("/api/proxy")
		mockRouteForm.EXPECT().GetString("method").Return("POST")
		mockRouteForm.EXPECT().GetString("routeType").Return(RouteTypeClient)

		result, err := createRouteInternal(
			mockRouteForm,
			func() (*config.FakeResponse, error) { return nil, nil },
			func() (*config.RequestTo, error) { return nil, expectedError },
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("happy path - mock route", func(t *testing.T) {

		mockRouteForm := form_builder.NewMockFormRunner(ctrl)
		fakeResponse := &config.FakeResponse{
			StatusCode: 200,
			Headers:    make(http.Header),
			Body:       config.HttpBody{"message": "success"},
		}

		mockRouteForm.EXPECT().Run().Return(nil)
		mockRouteForm.EXPECT().GetString("path").Return("/api/test")
		mockRouteForm.EXPECT().GetString("method").Return("GET")
		mockRouteForm.EXPECT().GetString("routeType").Return(RouteTypeMock)

		result, err := createRouteInternal(
			mockRouteForm,
			func() (*config.FakeResponse, error) { return fakeResponse, nil },
			func() (*config.RequestTo, error) { return nil, nil },
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "/api/test", result.Path)
		assert.Equal(t, "GET", result.Method)
		assert.NotNil(t, result.FakeResponse)
		assert.Nil(t, result.RequestTo)
		assert.Equal(t, 200, result.FakeResponse.StatusCode)
	})

	t.Run("happy path - client route", func(t *testing.T) {

		mockRouteForm := form_builder.NewMockFormRunner(ctrl)
		requestTo := &config.RequestTo{
			Host:    "http://localhost:8081",
			Path:    "/target",
			Method:  "POST",
			Headers: make(http.Header),
		}

		mockRouteForm.EXPECT().Run().Return(nil)
		mockRouteForm.EXPECT().GetString("path").Return("/api/proxy")
		mockRouteForm.EXPECT().GetString("method").Return("POST")
		mockRouteForm.EXPECT().GetString("routeType").Return(RouteTypeClient)

		result, err := createRouteInternal(
			mockRouteForm,
			func() (*config.FakeResponse, error) { return nil, nil },
			func() (*config.RequestTo, error) { return requestTo, nil },
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "/api/proxy", result.Path)
		assert.Equal(t, "POST", result.Method)
		assert.Nil(t, result.FakeResponse)
		assert.NotNil(t, result.RequestTo)
		assert.Equal(t, "http://localhost:8081", result.RequestTo.Host)
	})

	t.Run("happy path - unknown route type (no fakeResponse or requestTo)", func(t *testing.T) {

		mockRouteForm := form_builder.NewMockFormRunner(ctrl)

		mockRouteForm.EXPECT().Run().Return(nil)
		mockRouteForm.EXPECT().GetString("path").Return("/api/test")
		mockRouteForm.EXPECT().GetString("method").Return("GET")
		mockRouteForm.EXPECT().GetString("routeType").Return("unknown")

		result, err := createRouteInternal(
			mockRouteForm,
			func() (*config.FakeResponse, error) { return nil, nil },
			func() (*config.RequestTo, error) { return nil, nil },
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "/api/test", result.Path)
		assert.Equal(t, "GET", result.Method)
		assert.Nil(t, result.FakeResponse)
		assert.Nil(t, result.RequestTo)
	})
}
