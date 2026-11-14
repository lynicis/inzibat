package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/lynicis/inzibat/cmd/formbuilder"
	"github.com/lynicis/inzibat/config"
	_ "github.com/lynicis/inzibat/log"
)

var (
	httpMethods = []huh.Option[string]{
		{Key: "GET", Value: "GET"},
		{Key: "POST", Value: "POST"},
		{Key: "PUT", Value: "PUT"},
		{Key: "PATCH", Value: "PATCH"},
		{Key: "DELETE", Value: "DELETE"},
	}
	routeTypes = []huh.Option[string]{
		{Key: "Mock", Value: "mock"},
		{Key: "Client (Proxy)", Value: "client"},
	}
	createConfigFile           string
	createIsGlobalConfig       bool
	defaultFailureThreshold    = "5"
	defaultMinimumRequests     = "10"
	defaultOpenTimeoutMs       = "30000"
	defaultHalfOpenMaxRequests = "2"
	defaultSuccessThreshold    = "2"
)

func createRouteForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("path").
				Title("Route Path").
				Placeholder("/users").
				Validate(formbuilder.ValidatePath),
			huh.NewSelect[string]().
				Key("method").
				Title("HTTP Method").
				Options(httpMethods...),
			huh.NewSelect[string]().
				Key("routeType").
				Title("Route Type").
				Options(routeTypes...),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))
}

func createMockResponseFormInternal(
	statusFormRunner formbuilder.FormRunner,
	headersCollector func() (http.Header, error),
	bodyTypeFormRunner formbuilder.FormRunner,
	bodyCollector func() (config.HttpBody, error),
	bodyStringCollector func() (string, error),
) (*config.FakeResponse, error) {
	if err := statusFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to get status code: %w", err)
	}

	statusCodeStr := statusFormRunner.GetString("statusCode")
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status code %q: %w", statusCodeStr, err)
	}

	headers, err := headersCollector()
	if err != nil {
		return nil, fmt.Errorf("failed to collect headers: %w", err)
	}

	if err := bodyTypeFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to select body type: %w", err)
	}

	bodyType := bodyTypeFormRunner.GetString("bodyType")

	fakeResponse := &config.FakeResponse{
		StatusCode: statusCode,
		Headers:    headers,
	}

	switch bodyType {
	case BodyTypeBody:
		body, err := bodyCollector()
		if err != nil {
			return nil, fmt.Errorf("failed to collect body: %w", err)
		}
		fakeResponse.Body = body
	case BodyTypeBodyString:
		bodyString, err := bodyStringCollector()
		if err != nil {
			return nil, fmt.Errorf("failed to collect body string: %w", err)
		}
		fakeResponse.BodyString = bodyString
	}

	return fakeResponse, nil
}

func createMockResponseForm() (*config.FakeResponse, error) {
	status := strconv.Itoa(http.StatusOK)
	statusForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("statusCode").
				Title("Status Code").
				Placeholder(status).
				Value(&status).
				Validate(formbuilder.ValidateStatusCode),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

	bodyForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("bodyType").
				Title("Body Type").
				Options([]huh.Option[string]{
					{Key: "Body (JSON object)", Value: BodyTypeBody},
					{Key: "BodyString (string)", Value: BodyTypeBodyString},
					{Key: "Skip", Value: formbuilder.SourceSkip},
				}...),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

	statusFormRunner := &formbuilder.HuhFormRunner{Form: statusForm}
	bodyTypeFormRunner := &formbuilder.HuhFormRunner{Form: bodyForm}

	return createMockResponseFormInternal(
		statusFormRunner,
		formbuilder.CollectHeaders,
		bodyTypeFormRunner,
		formbuilder.CollectBody,
		formbuilder.CollectBodyString,
	)
}

func createClientRequestFormInternal(
	basicFormRunner formbuilder.FormRunner,
	headersCollector func() (http.Header, error),
	bodyTypeFormRunner formbuilder.FormRunner,
	bodyCollector func() (config.HttpBody, error),
	optionsFormRunner formbuilder.FormRunner,
) (*config.RequestTo, error) {
	if err := basicFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to get basic request info: %w", err)
	}

	host := basicFormRunner.GetString("host")
	targetPath := basicFormRunner.GetString("path")
	targetMethod := basicFormRunner.GetString("method")

	headers, err := headersCollector()
	if err != nil {
		return nil, fmt.Errorf("failed to collect headers: %w", err)
	}

	if err := bodyTypeFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to select body type: %w", err)
	}

	body, err := collectClientRequestBody(bodyTypeFormRunner, bodyCollector)
	if err != nil {
		return nil, err
	}

	if err := optionsFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	circuitBreakerConfig, err := collectCircuitBreakerConfig(optionsFormRunner)
	if err != nil {
		return nil, err
	}

	return &config.RequestTo{
		Host:                   host,
		Path:                   targetPath,
		Method:                 targetMethod,
		Headers:                headers,
		Body:                   body,
		PassWithRequestBody:    optionsFormRunner.GetBool("passWithRequestBody"),
		PassWithRequestHeaders: optionsFormRunner.GetBool("passWithRequestHeaders"),
		InErrorReturn500:       optionsFormRunner.GetBool("inErrorReturn500"),
		CircuitBreaker:         circuitBreakerConfig,
	}, nil
}

func collectClientRequestBody(
	bodyTypeFormRunner formbuilder.FormRunner,
	bodyCollector func() (config.HttpBody, error),
) (config.HttpBody, error) {
	if bodyTypeFormRunner.GetString("bodyType") != BodyTypeStructured {
		return nil, nil
	}

	body, err := bodyCollector()
	if err != nil {
		return nil, fmt.Errorf("failed to collect body: %w", err)
	}

	return body, nil
}

func collectCircuitBreakerConfig(
	optionsFormRunner formbuilder.FormRunner,
) (*config.CircuitBreakerConfig, error) {
	circuitBreakerEnabled := optionsFormRunner.GetBool("circuitBreakerEnabled")
	circuitBreakerConfig := &config.CircuitBreakerConfig{
		Enabled: config.BoolPointer(circuitBreakerEnabled),
	}

	if !circuitBreakerEnabled {
		return circuitBreakerConfig, nil
	}

	failureThreshold, err := parsePositiveInt(
		optionsFormRunner.GetString("failureThreshold"),
		"failureThreshold",
	)
	if err != nil {
		return nil, err
	}

	minimumRequests, err := parsePositiveInt(
		optionsFormRunner.GetString("minimumRequests"),
		"minimumRequests",
	)
	if err != nil {
		return nil, err
	}

	openTimeoutMs, err := parsePositiveInt(
		optionsFormRunner.GetString("openTimeoutMs"),
		"openTimeoutMs",
	)
	if err != nil {
		return nil, err
	}

	halfOpenMaxRequests, err := parsePositiveInt(
		optionsFormRunner.GetString("halfOpenMaxRequests"),
		"halfOpenMaxRequests",
	)
	if err != nil {
		return nil, err
	}

	successThreshold, err := parsePositiveInt(
		optionsFormRunner.GetString("successThreshold"),
		"successThreshold",
	)
	if err != nil {
		return nil, err
	}

	circuitBreakerConfig.FailureThreshold = failureThreshold
	circuitBreakerConfig.MinimumRequests = minimumRequests
	circuitBreakerConfig.OpenTimeoutMs = openTimeoutMs
	circuitBreakerConfig.HalfOpenMaxRequests = halfOpenMaxRequests
	circuitBreakerConfig.SuccessThreshold = successThreshold

	return circuitBreakerConfig, nil
}

func parsePositiveInt(value string, key string) (int, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}

	if parsedValue <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return parsedValue, nil
}

func createClientRequestForm() (*config.RequestTo, error) {
	basicForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("host").
				Title("Target Host URL").
				Placeholder("http://localhost:8081").
				Validate(formbuilder.ValidateHost),
			huh.NewInput().
				Key("path").
				Title("Target Path").
				Placeholder("/api/users").
				Validate(formbuilder.ValidateHost),
			huh.NewSelect[string]().
				Key("method").
				Title("Target HTTP Method").
				Options(httpMethods...),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

	bodyForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("bodyType").
				Title("Request Body Type").
				Options([]huh.Option[string]{
					{Key: "Structured (JSON object)", Value: BodyTypeStructured},
					{Key: "Skip", Value: formbuilder.SourceSkip},
				}...),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

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
			huh.NewConfirm().
				Key("circuitBreakerEnabled").
				Title("Enable Circuit Breaker"),
			huh.NewInput().
				Key("failureThreshold").
				Title("Circuit Breaker Failure Threshold").
				Value(&defaultFailureThreshold).
				Validate(formbuilder.ValidatePositiveInt),
			huh.NewInput().
				Key("minimumRequests").
				Title("Circuit Breaker Minimum Requests").
				Value(&defaultMinimumRequests).
				Validate(formbuilder.ValidatePositiveInt),
			huh.NewInput().
				Key("openTimeoutMs").
				Title("Circuit Breaker Open Timeout (ms)").
				Value(&defaultOpenTimeoutMs).
				Validate(formbuilder.ValidatePositiveInt),
			huh.NewInput().
				Key("halfOpenMaxRequests").
				Title("Circuit Breaker Half-Open Max Requests").
				Value(&defaultHalfOpenMaxRequests).
				Validate(formbuilder.ValidatePositiveInt),
			huh.NewInput().
				Key("successThreshold").
				Title("Circuit Breaker Success Threshold").
				Value(&defaultSuccessThreshold).
				Validate(formbuilder.ValidatePositiveInt),
		),
	).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithProgramOptions(tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout))

	basicFormRunner := &formbuilder.HuhFormRunner{Form: basicForm}
	bodyTypeFormRunner := &formbuilder.HuhFormRunner{Form: bodyForm}
	optionsFormRunner := &formbuilder.HuhFormRunner{Form: optionsForm}

	return createClientRequestFormInternal(
		basicFormRunner,
		formbuilder.CollectHeaders,
		bodyTypeFormRunner,
		formbuilder.CollectBody,
		optionsFormRunner,
	)
}

func createRouteInternal(
	routeFormRunner formbuilder.FormRunner,
	mockResponseFormCreator func() (*config.FakeResponse, error),
	clientRequestFormCreator func() (*config.RequestTo, error),
) (*config.Route, error) {
	if err := routeFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	path := routeFormRunner.GetString("path")
	method := routeFormRunner.GetString("method")
	routeType := routeFormRunner.GetString("routeType")

	var fakeResponse *config.FakeResponse
	var requestTo *config.RequestTo

	switch routeType {
	case RouteTypeMock:
		var err error
		fakeResponse, err = mockResponseFormCreator()
		if err != nil {
			return nil, err
		}
	case RouteTypeClient:
		var err error
		requestTo, err = clientRequestFormCreator()
		if err != nil {
			return nil, err
		}
	}

	return &config.Route{
		Method:       method,
		Path:         path,
		FakeResponse: fakeResponse,
		RequestTo:    requestTo,
	}, nil
}

func createRoute() (*config.Route, error) {
	routeForm := createRouteForm()
	routeFormRunner := &formbuilder.HuhFormRunner{Form: routeForm}

	return createRouteInternal(
		routeFormRunner,
		createMockResponseForm,
		createClientRequestForm,
	)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"create-route", "c"},
	Short:   "Create a new route",
	Long: `Create a new route interactively using a form-based CLI.

The route will be added to the configuration file (in order of precedence):
  1. The file specified by the --config flag
  2. The file specified by the INZIBAT_CONFIG_FILE environment variable
  3. inzibat.json in the current working directory
  4. ~/.inzibat.config.json if --global flag is used`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		route, err := createRoute()
		if err != nil {
			zap.L().Fatal("failed to create route", zap.Error(err))
		}

		cfgFilePath, err := config.ResolveConfigPath(createIsGlobalConfig, createConfigFile)
		if err != nil {
			zap.L().Fatal("failed to resolve config file path", zap.Error(err))
		}

		cfg, err := config.ReadOrCreateConfig(cfgFilePath)
		if err != nil {
			zap.L().Fatal("failed to read/create config", zap.Error(err))
		}

		cfg.Routes = append(cfg.Routes, *route)

		if err := config.WriteConfig(cfg, cfgFilePath); err != nil {
			zap.L().Fatal("failed to write config", zap.Error(err))
		}

		zap.L().Info("Route created successfully", zap.String("config_file", cfgFilePath))
	},
}

func init() {
	createCmd.Flags().StringVarP(
		&createConfigFile,
		"config",
		"c",
		"",
		"Path to the configuration file",
	)
	createCmd.Flags().BoolVarP(
		&createIsGlobalConfig,
		"global",
		"g",
		false,
		"Use the global config file (~/.inzibat.config.json)",
	)
	rootCmd.AddCommand(createCmd)
}
