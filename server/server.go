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

	"inzibat/client/http"
	"inzibat/config"
	"inzibat/handler"
	_ "inzibat/log"
	"inzibat/router"
)

func StartServer(configFile string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return StartServerWithContext(ctx, configFile)
}

func StartServerWithContext(ctx context.Context, configFile string) error {
	if configFile != "" {
		absPath, err := config.ResolveAbsolutePath(configFile)
		if err != nil {
			return fmt.Errorf("failed to resolve config file path: %w", err)
		}
		configFile = absPath

		originalEnv := os.Getenv(config.EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				if err := os.Setenv(config.EnvironmentVariableConfigFileName, originalEnv); err != nil {
					zap.L().Warn("failed to restore environment variable", zap.Error(err))
				}
			} else {
				if err := os.Unsetenv(config.EnvironmentVariableConfigFileName); err != nil {
					zap.L().Warn("failed to unset environment variable", zap.Error(err))
				}
			}
		}()

		if err := os.Setenv(config.EnvironmentVariableConfigFileName, configFile); err != nil {
			return fmt.Errorf("failed to set environment variable: %w", err)
		}
	}

	validator := validatorPkg.New()
	configLoader := config.NewLoader(validator, false)
	cfg, err := configLoader.Read()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

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

	go func() {
		if err = fiberApp.Listen(cfg.GetServerAddr()); err != nil {
			zap.L().Fatal("failed to start http server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	if err = fiberApp.ShutdownWithTimeout(5 * time.Second); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}

	return nil
}
