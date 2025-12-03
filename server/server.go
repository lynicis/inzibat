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
	"github.com/lynicis/inzibat/router"
)

func StartServer(configFile string, isGlobalConfig bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return StartServerWithContext(ctx, configFile, isGlobalConfig)
}

func StartServerWithContext(ctx context.Context, configFile string, isGlobalConfig bool) error {
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

	fiberApp, err := setupServer(cfg)
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

func setupServer(cfg *config.Cfg) (*fiber.App, error) {
	endpointHandler := &handler.EndpointHandler{
		RouteConfig: &cfg.Routes,
	}
	httpClient := http.NewHttpClient()
	clientHandler := &handler.ClientHandler{
		Client:      httpClient,
		RouteConfig: &cfg.Routes,
	}

	fiberApp := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
	})
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
