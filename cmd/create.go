package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"inzibat/cmd/form_builder"
	"inzibat/config"
	_ "inzibat/log"
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
)

func createRouteForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("path").
				Title("Route Path").
				Placeholder("/users").
				Validate(form_builder.ValidatePath),
			huh.NewSelect[string]().
				Key("method").
				Title("HTTP Method").
				Options(httpMethods...),
			huh.NewSelect[string]().
				Key("routeType").
				Title("Route Type").
				Options(routeTypes...),
		),
	)
}

func createMockResponseFormInternal(
	statusFormRunner form_builder.FormRunner,
	headersCollector func() (http.Header, error),
	bodyTypeFormRunner form_builder.FormRunner,
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
				Validate(form_builder.ValidateStatusCode),
		),
	)

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

	statusFormRunner := &form_builder.HuhFormRunner{Form: statusForm}
	bodyTypeFormRunner := &form_builder.HuhFormRunner{Form: bodyForm}

	return createMockResponseFormInternal(
		statusFormRunner,
		form_builder.CollectHeaders,
		bodyTypeFormRunner,
		form_builder.CollectBody,
		form_builder.CollectBodyString,
	)
}

func createClientRequestFormInternal(
	basicFormRunner form_builder.FormRunner,
	headersCollector func() (http.Header, error),
	bodyTypeFormRunner form_builder.FormRunner,
	bodyCollector func() (config.HttpBody, error),
	optionsFormRunner form_builder.FormRunner,
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

	var body config.HttpBody
	bodyType := bodyTypeFormRunner.GetString("bodyType")
	if bodyType == BodyTypeStructured {
		body, err = bodyCollector()
		if err != nil {
			return nil, fmt.Errorf("failed to collect body: %w", err)
		}
	}

	if err := optionsFormRunner.Run(); err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
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
	}, nil
}

func createClientRequestForm() (*config.RequestTo, error) {
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

	basicFormRunner := &form_builder.HuhFormRunner{Form: basicForm}
	bodyTypeFormRunner := &form_builder.HuhFormRunner{Form: bodyForm}
	optionsFormRunner := &form_builder.HuhFormRunner{Form: optionsForm}

	return createClientRequestFormInternal(
		basicFormRunner,
		form_builder.CollectHeaders,
		bodyTypeFormRunner,
		form_builder.CollectBody,
		optionsFormRunner,
	)
}

func createRouteInternal(
	routeFormRunner form_builder.FormRunner,
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
	routeFormRunner := &form_builder.HuhFormRunner{Form: routeForm}

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
	Long:    `Create a new route interactively using a form-based CLI.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		route, err := createRoute()
		if err != nil {
			zap.L().Fatal("failed to create route", zap.Error(err))
		}

		cfgLoader := config.NewLoader(nil, false)
		cfgFilePath := cfgLoader.Filepath

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
	rootCmd.AddCommand(createCmd)
}
