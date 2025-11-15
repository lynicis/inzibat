package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	if err := statusForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to get status code: %w", err)
	}

	statusCodeStr := statusForm.GetString("statusCode")
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status code %q: %w", statusCodeStr, err)
	}

	headers, err := form_builder.CollectHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to collect headers: %w", err)
	}

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

	if err := bodyForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to select body type: %w", err)
	}

	bodyType := bodyForm.GetString("bodyType")

	fakeResponse := config.FakeResponse{
		StatusCode: statusCode,
		Headers:    headers,
	}

	switch bodyType {
	case BodyTypeBody:
		body, err := form_builder.CollectBody()
		if err != nil {
			return nil, fmt.Errorf("failed to collect body: %w", err)
		}
		fakeResponse.Body = body
	case BodyTypeBodyString:
		bodyString, err := form_builder.CollectBodyString()
		if err != nil {
			return nil, fmt.Errorf("failed to collect body string: %w", err)
		}
		fakeResponse.BodyString = bodyString
	}

	return &fakeResponse, nil
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

	if err := basicForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to get basic request info: %w", err)
	}

	host := basicForm.GetString("host")
	targetPath := basicForm.GetString("path")
	targetMethod := basicForm.GetString("method")

	headers, err := form_builder.CollectHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed to collect headers: %w", err)
	}

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

	if err := bodyForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to select body type: %w", err)
	}

	var body config.HttpBody
	bodyType := bodyForm.GetString("bodyType")
	if bodyType == BodyTypeStructured {
		body, err = form_builder.CollectBody()
		if err != nil {
			return nil, fmt.Errorf("failed to collect body: %w", err)
		}
	}

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

	if err := optionsForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	return &config.RequestTo{
		Host:                   host,
		Path:                   targetPath,
		Method:                 targetMethod,
		Headers:                headers,
		Body:                   body,
		PassWithRequestBody:    optionsForm.GetBool("passWithRequestBody"),
		PassWithRequestHeaders: optionsForm.GetBool("passWithRequestHeaders"),
		InErrorReturn500:       optionsForm.GetBool("inErrorReturn500"),
	}, nil
}

func createRoute() (*config.Route, error) {
	routeForm := createRouteForm()
	if err := routeForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	path := routeForm.GetString("path")
	method := routeForm.GetString("method")
	routeType := routeForm.GetString("routeType")

	var fakeResponse config.FakeResponse
	var requestTo config.RequestTo

	switch routeType {
	case RouteTypeMock:
		fakeResponsePtr, err := createMockResponseForm()
		if err != nil {
			return nil, err
		}
		fakeResponse = *fakeResponsePtr
	case RouteTypeClient:
		requestToPtr, err := createClientRequestForm()
		if err != nil {
			return nil, err
		}
		requestTo = *requestToPtr
	}

	return &config.Route{
		Method:       method,
		Path:         path,
		FakeResponse: fakeResponse,
		RequestTo:    requestTo,
	}, nil
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

		homeDir, err := os.UserHomeDir()
		if err != nil {
			zap.L().Fatal("failed to get current user's home directory", zap.Error(err))
		}

		globalConfigFilePath := filepath.Join(homeDir, config.DefaultConfigFileName)
		cfg, err := config.ReadOrCreateConfig(globalConfigFilePath)
		if err != nil {
			zap.L().Fatal("failed to read/create config", zap.Error(err))
		}

		cfg.Routes = append(cfg.Routes, *route)

		if err := config.WriteConfig(cfg, globalConfigFilePath); err != nil {
			zap.L().Fatal("failed to write config", zap.Error(err))
		}

		zap.L().Info("Route created successfully", zap.String("config_file", globalConfigFilePath))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
