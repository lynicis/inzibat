package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lynicis/inzibat/cmd/formbuilder"
	"github.com/lynicis/inzibat/config"
)

func TestCreateRouteFormHeadless(t *testing.T) {
	var runner *formbuilder.HuhFormRunner

	withHeadlessTerminal(t, func() {
		runner = &formbuilder.HuhFormRunner{Form: createRouteForm()}
		require.NoError(t, runner.Run())
	},
		"/users\r",
		formbuilderHeadlessSelect(0), // method: GET
		formbuilderHeadlessSelect(0), // route type: mock
	)

	assert.Equal(t, "/users", runner.GetString("path"))
	assert.Equal(t, "GET", runner.GetString("method"))
	assert.Equal(t, RouteTypeMock, runner.GetString("routeType"))
}

func TestCreateMockResponseFormHeadlessBodyString(t *testing.T) {
	var fakeResponse *config.FakeResponse

	withHeadlessTerminal(t, func() {
		var err error
		fakeResponse, err = createMockResponseForm()
		require.NoError(t, err)
	},
		formbuilderHeadlessEnter,     // status code: default 200
		formbuilderHeadlessSelect(2), // header source: skip
		formbuilderHeadlessSelect(1), // body type: body string
		formbuilderHeadlessSelect(1), // body string source: form
		"created via form\r",
	)

	require.NotNil(t, fakeResponse)
	assert.Equal(t, 200, fakeResponse.StatusCode)
	assert.Empty(t, fakeResponse.Headers)
	assert.Equal(t, "created via form", fakeResponse.BodyString)
	assert.Nil(t, fakeResponse.Body)
}

func TestCreateMockResponseFormHeadlessBodyObject(t *testing.T) {
	var fakeResponse *config.FakeResponse

	withHeadlessTerminal(t, func() {
		var err error
		fakeResponse, err = createMockResponseForm()
		require.NoError(t, err)
	},
		formbuilderHeadlessEnter,     // status code: default 200
		formbuilderHeadlessSelect(2), // header source: skip
		formbuilderHeadlessSelect(0), // body type: json object
		formbuilderHeadlessSelect(1), // body source: form
		"message\r",
		"success\r",
		formbuilderHeadlessNo,
	)

	require.NotNil(t, fakeResponse)
	assert.Equal(t, 200, fakeResponse.StatusCode)
	assert.Equal(t, config.HttpBody{"message": "success"}, fakeResponse.Body)
	assert.Empty(t, fakeResponse.BodyString)
}

func TestCreateMockResponseFormHeadlessHeadersFromFile(t *testing.T) {
	headerFile := writeTempFile(t, "headers.json", `{"X-File": "loaded"}`)

	var fakeResponse *config.FakeResponse

	withHeadlessTerminal(t, func() {
		var err error
		fakeResponse, err = createMockResponseForm()
		require.NoError(t, err)
	},
		formbuilderHeadlessEnter,     // status code: default 200
		formbuilderHeadlessSelect(0), // header source: file
		headerFile+formbuilderHeadlessEnter,
		formbuilderHeadlessSelect(2), // body type: skip
	)

	require.NotNil(t, fakeResponse)
	assert.Equal(t, "loaded", fakeResponse.Headers.Get("X-File"))
	assert.Nil(t, fakeResponse.Body)
}

func clientRequestSegments(enableBreaker string) []string {
	return append([]string{
		"http://localhost:8081\r",
		"/upstream\r",
		formbuilderHeadlessSelect(0), // method: GET
		formbuilderHeadlessSelect(2), // header source: skip
		formbuilderHeadlessSelect(1), // body type: skip
		formbuilderHeadlessNo,        // pass with request body
		formbuilderHeadlessNo,        // pass with request headers
		formbuilderHeadlessNo,        // return 500 on error
		enableBreaker,
	}, breakerDefaults()...)
}

// breakerDefaults accepts the five circuit-breaker inputs' preset values.
func breakerDefaults() []string {
	return []string{
		formbuilderHeadlessEnter,
		formbuilderHeadlessEnter,
		formbuilderHeadlessEnter,
		formbuilderHeadlessEnter,
		formbuilderHeadlessEnter,
	}
}

func TestCreateClientRequestFormHeadless(t *testing.T) {
	var requestTo *config.RequestTo

	withHeadlessTerminal(t, func() {
		var err error
		requestTo, err = createClientRequestForm()
		require.NoError(t, err)
	}, clientRequestSegments(formbuilderHeadlessYes)...)

	require.NotNil(t, requestTo)
	assert.Equal(t, "http://localhost:8081", requestTo.Host)
	assert.Equal(t, "/upstream", requestTo.Path)
	assert.Equal(t, "GET", requestTo.Method)
	assert.False(t, requestTo.PassWithRequestBody)
	assert.False(t, requestTo.PassWithRequestHeaders)
	assert.False(t, requestTo.InErrorReturn500)
	require.NotNil(t, requestTo.CircuitBreaker)
	require.NotNil(t, requestTo.CircuitBreaker.Enabled)
	assert.True(t, *requestTo.CircuitBreaker.Enabled)
	assert.Equal(t, 5, requestTo.CircuitBreaker.FailureThreshold)
	assert.Equal(t, 10, requestTo.CircuitBreaker.MinimumRequests)
	assert.Equal(t, 30000, requestTo.CircuitBreaker.OpenTimeoutMs)
	assert.Equal(t, 2, requestTo.CircuitBreaker.HalfOpenMaxRequests)
	assert.Equal(t, 2, requestTo.CircuitBreaker.SuccessThreshold)
}

func TestCreateClientRequestFormHeadlessBreakerDisabled(t *testing.T) {
	var requestTo *config.RequestTo

	withHeadlessTerminal(t, func() {
		var err error
		requestTo, err = createClientRequestForm()
		require.NoError(t, err)
	}, clientRequestSegments(formbuilderHeadlessNo)...)

	require.NotNil(t, requestTo)
	require.NotNil(t, requestTo.CircuitBreaker)
	require.NotNil(t, requestTo.CircuitBreaker.Enabled)
	assert.False(t, *requestTo.CircuitBreaker.Enabled)
	assert.Zero(t, requestTo.CircuitBreaker.FailureThreshold)
}

func TestCreateRouteHeadlessMock(t *testing.T) {
	var route *config.Route

	withHeadlessTerminal(t, func() {
		var err error
		route, err = createRoute()
		require.NoError(t, err)
	},
		"/users\r",
		formbuilderHeadlessSelect(0), // method: GET
		formbuilderHeadlessSelect(0), // route type: mock
		formbuilderHeadlessEnter,     // status code: default 200
		formbuilderHeadlessSelect(2), // header source: skip
		formbuilderHeadlessSelect(1), // body type: body string
		formbuilderHeadlessSelect(1), // body string source: form
		"created via form\r",
	)

	require.NotNil(t, route)
	assert.Equal(t, "GET", route.Method)
	assert.Equal(t, "/users", route.Path)
	require.NotNil(t, route.FakeResponse)
	assert.Equal(t, 200, route.FakeResponse.StatusCode)
	assert.Equal(t, "created via form", route.FakeResponse.BodyString)
	assert.Nil(t, route.RequestTo)
}

func TestCreateRouteHeadlessClient(t *testing.T) {
	var route *config.Route

	withHeadlessTerminal(t, func() {
		var err error
		route, err = createRoute()
		require.NoError(t, err)
	}, append([]string{
		"/proxy\r",
		formbuilderHeadlessSelect(0), // method: GET
		formbuilderHeadlessSelect(1), // route type: client
	}, clientRequestSegments(formbuilderHeadlessNo)...)...)

	require.NotNil(t, route)
	assert.Equal(t, "GET", route.Method)
	assert.Equal(t, "/proxy", route.Path)
	assert.Nil(t, route.FakeResponse)
	require.NotNil(t, route.RequestTo)
	assert.Equal(t, "http://localhost:8081", route.RequestTo.Host)
	assert.Equal(t, "/upstream", route.RequestTo.Path)
}
