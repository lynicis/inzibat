package cmd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/lynicis/inzibat/cmd/formbuilder"
	"github.com/lynicis/inzibat/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateMockResponseFormInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - statusFormRunner.Run() returns error", func(t *testing.T) {

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)

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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		headers := make(http.Header)

		mockStatusForm.EXPECT().Run().Return(nil)
		mockStatusForm.EXPECT().GetString("statusCode").Return("204")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(formbuilder.SourceSkip)

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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockStatusForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
		expectedError := fmt.Errorf("options form error")

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(formbuilder.SourceSkip)
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

	t.Run("error path - invalid circuit breaker threshold", func(t *testing.T) {
		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(formbuilder.SourceSkip)
		mockOptionsForm.EXPECT().Run().Return(nil)
		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("abc")

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return make(http.Header), nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse failureThreshold")
	})

	t.Run("happy path - with structured body", func(t *testing.T) {

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
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
		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(false)

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
		assert.NotNil(t, result.CircuitBreaker)
		assert.NotNil(t, result.CircuitBreaker.Enabled)
		assert.False(t, *result.CircuitBreaker.Enabled)
	})

	t.Run("happy path - skip body", func(t *testing.T) {

		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
		headers := make(http.Header)

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(formbuilder.SourceSkip)
		mockOptionsForm.EXPECT().Run().Return(nil)
		mockOptionsForm.EXPECT().GetBool("passWithRequestBody").Return(false)
		mockOptionsForm.EXPECT().GetBool("passWithRequestHeaders").Return(false)
		mockOptionsForm.EXPECT().GetBool("inErrorReturn500").Return(true)
		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(false)

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
		assert.NotNil(t, result.CircuitBreaker)
		assert.NotNil(t, result.CircuitBreaker.Enabled)
		assert.False(t, *result.CircuitBreaker.Enabled)
	})

	t.Run("happy path - circuit breaker enabled", func(t *testing.T) {
		mockBasicForm := formbuilder.NewMockFormRunner(ctrl)
		mockBodyTypeForm := formbuilder.NewMockFormRunner(ctrl)
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)
		headers := make(http.Header)

		mockBasicForm.EXPECT().Run().Return(nil)
		mockBasicForm.EXPECT().GetString("host").Return("http://localhost:8081")
		mockBasicForm.EXPECT().GetString("path").Return("/api/users")
		mockBasicForm.EXPECT().GetString("method").Return("GET")
		mockBodyTypeForm.EXPECT().Run().Return(nil)
		mockBodyTypeForm.EXPECT().GetString("bodyType").Return(formbuilder.SourceSkip)
		mockOptionsForm.EXPECT().Run().Return(nil)
		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("6")
		mockOptionsForm.EXPECT().GetString("minimumRequests").Return("11")
		mockOptionsForm.EXPECT().GetString("openTimeoutMs").Return("20000")
		mockOptionsForm.EXPECT().GetString("halfOpenMaxRequests").Return("3")
		mockOptionsForm.EXPECT().GetString("successThreshold").Return("2")
		mockOptionsForm.EXPECT().GetBool("passWithRequestBody").Return(false)
		mockOptionsForm.EXPECT().GetBool("passWithRequestHeaders").Return(false)
		mockOptionsForm.EXPECT().GetBool("inErrorReturn500").Return(false)

		result, err := createClientRequestFormInternal(
			mockBasicForm,
			func() (http.Header, error) { return headers, nil },
			mockBodyTypeForm,
			func() (config.HttpBody, error) { return nil, nil },
			mockOptionsForm,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.CircuitBreaker)
		assert.NotNil(t, result.CircuitBreaker.Enabled)
		assert.True(t, *result.CircuitBreaker.Enabled)
		assert.Equal(t, 6, result.CircuitBreaker.FailureThreshold)
		assert.Equal(t, 11, result.CircuitBreaker.MinimumRequests)
		assert.Equal(t, 20000, result.CircuitBreaker.OpenTimeoutMs)
		assert.Equal(t, 3, result.CircuitBreaker.HalfOpenMaxRequests)
		assert.Equal(t, 2, result.CircuitBreaker.SuccessThreshold)
	})
}

func TestParsePositiveInt(t *testing.T) {
	t.Run("error path - zero value", func(t *testing.T) {
		result, err := parsePositiveInt("0", "testField")

		assert.Error(t, err)
		assert.Equal(t, 0, result)
		assert.Contains(t, err.Error(), "testField must be greater than zero")
	})

	t.Run("error path - negative value", func(t *testing.T) {
		result, err := parsePositiveInt("-5", "testField")

		assert.Error(t, err)
		assert.Equal(t, 0, result)
		assert.Contains(t, err.Error(), "testField must be greater than zero")
	})

	t.Run("happy path - positive value", func(t *testing.T) {
		result, err := parsePositiveInt("10", "testField")

		assert.NoError(t, err)
		assert.Equal(t, 10, result)
	})
}

func TestCollectCircuitBreakerConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - invalid minimumRequests", func(t *testing.T) {
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)

		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("5")
		mockOptionsForm.EXPECT().GetString("minimumRequests").Return("invalid")

		result, err := collectCircuitBreakerConfig(mockOptionsForm)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse minimumRequests")
	})

	t.Run("error path - invalid openTimeoutMs", func(t *testing.T) {
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)

		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("5")
		mockOptionsForm.EXPECT().GetString("minimumRequests").Return("10")
		mockOptionsForm.EXPECT().GetString("openTimeoutMs").Return("invalid")

		result, err := collectCircuitBreakerConfig(mockOptionsForm)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse openTimeoutMs")
	})

	t.Run("error path - invalid halfOpenMaxRequests", func(t *testing.T) {
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)

		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("5")
		mockOptionsForm.EXPECT().GetString("minimumRequests").Return("10")
		mockOptionsForm.EXPECT().GetString("openTimeoutMs").Return("30000")
		mockOptionsForm.EXPECT().GetString("halfOpenMaxRequests").Return("invalid")

		result, err := collectCircuitBreakerConfig(mockOptionsForm)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse halfOpenMaxRequests")
	})

	t.Run("error path - invalid successThreshold", func(t *testing.T) {
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)

		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("5")
		mockOptionsForm.EXPECT().GetString("minimumRequests").Return("10")
		mockOptionsForm.EXPECT().GetString("openTimeoutMs").Return("30000")
		mockOptionsForm.EXPECT().GetString("halfOpenMaxRequests").Return("2")
		mockOptionsForm.EXPECT().GetString("successThreshold").Return("invalid")

		result, err := collectCircuitBreakerConfig(mockOptionsForm)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse successThreshold")
	})

	t.Run("error path - non-positive failureThreshold", func(t *testing.T) {
		mockOptionsForm := formbuilder.NewMockFormRunner(ctrl)

		mockOptionsForm.EXPECT().GetBool("circuitBreakerEnabled").Return(true)
		mockOptionsForm.EXPECT().GetString("failureThreshold").Return("0")

		result, err := collectCircuitBreakerConfig(mockOptionsForm)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failureThreshold must be greater than zero")
	})
}

func TestCreateRouteInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error path - routeFormRunner.Run() returns error", func(t *testing.T) {
		mockRouteForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockRouteForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockRouteForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockRouteForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockRouteForm := formbuilder.NewMockFormRunner(ctrl)
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

		mockRouteForm := formbuilder.NewMockFormRunner(ctrl)

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
