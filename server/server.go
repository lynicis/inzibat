package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/lynicis/inzibat/client/http"
	"github.com/lynicis/inzibat/config"
	"github.com/lynicis/inzibat/handler"
	_ "github.com/lynicis/inzibat/log"
	"github.com/lynicis/inzibat/recorder"
	"github.com/lynicis/inzibat/router"
)

func StartServer(configFile string, isGlobalConfig bool, recordEnabled bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return StartServerWithContext(ctx, configFile, isGlobalConfig, recordEnabled)
}

func StartServerWithContext(
	ctx context.Context,
	configFile string,
	isGlobalConfig bool,
	recordEnabled bool,
) error {
	var resolvedPath string
	if configFile != "" {
		absPath, err := config.ResolveAbsolutePath(configFile)
		if err != nil {
			return fmt.Errorf("failed to resolve config file path: %w", err)
		}
		resolvedPath = absPath
	}

	cfg, err := loadConfig(resolvedPath, isGlobalConfig)
	if err != nil {
		return err
	}

	fiberApp, err := setupServer(cfg, recordEnabled)
	if err != nil {
		return err
	}

	return runServer(ctx, fiberApp, cfg)
}

func loadConfig(explicitPath string, isGlobalConfig bool) (*config.Cfg, error) {
	validator := validatorPkg.New()
	configLoader := config.NewLoader(validator, isGlobalConfig, explicitPath)
	cfg, err := configLoader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return cfg, nil
}

func setupServer(cfg *config.Cfg, recordEnabled bool) (*fiber.App, error) {
	endpointHandler := &handler.EndpointHandler{
		RouteConfig: &cfg.Routes,
	}
	httpClient := http.NewHttpClient()
	circuitBreakerStore, err := handler.NewCircuitBreakerStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize circuit breaker store: %w", err)
	}

	circuitBreakerRouteKeys := make(map[int]string)
	for routeIndex := range cfg.Routes {
		route := cfg.Routes[routeIndex]
		if route.RequestTo == nil || route.RequestTo.CircuitBreaker == nil {
			continue
		}

		if route.RequestTo.CircuitBreaker.Enabled != nil && *route.RequestTo.CircuitBreaker.Enabled {
			routeKey := handler.BuildCircuitBreakerRouteKey(route)
			if err = circuitBreakerStore.Seed(routeKey, *route.RequestTo.CircuitBreaker); err != nil {
				return nil, fmt.Errorf("failed to seed circuit breaker store: %w", err)
			}

			circuitBreakerRouteKeys[routeIndex] = routeKey
		}
	}

	clientHandler := &handler.ClientHandler{
		Client:                  httpClient,
		RouteConfig:             &cfg.Routes,
		CircuitBreakerStore:     circuitBreakerStore,
		CircuitBreakerRouteKeys: circuitBreakerRouteKeys,
	}

	fiberApp := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
		ReadBufferSize:        4 * 1024 * 1024,
	})

	if recordEnabled {
		recordStore := recorder.NewStore(recorder.DefaultStoreCapacity)
		fiberApp.Use(recorder.NewRecorderMiddleware(recordStore))
		recorder.RegisterAdminRoutes(fiberApp, recordStore)
		zap.L().Info("🔴 Request recording enabled")
	}

	mainRouter := &router.MainRouter{
		Config:          cfg,
		FiberApp:        fiberApp,
		EndpointHandler: endpointHandler,
		ClientHandler:   clientHandler,
	}
	mainRouter.CreateRoutes()

	zap.L().Info("🫡 INZIBAT 🪖",
		zap.Int("open_routes", len(cfg.Routes)),
		zap.Int("server_port", cfg.ServerPort),
	)

	return fiberApp, nil
}

func runServer(ctx context.Context, fiberApp *fiber.App, cfg *config.Cfg) error {
	var serverErr error
	go func() {
		if err := fiberApp.Listen(cfg.GetServerAddr()); err != nil {
			zap.L().Fatal("failed to start http server", zap.Error(err))
			serverErr = err
		}
	}()

	<-ctx.Done()

	if err := fiberApp.ShutdownWithTimeout(5 * time.Second); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}

	return serverErr
}
