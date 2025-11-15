package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Reader struct {
	ConfigReader ReaderStrategy
	Validator    *validator.Validate
	Filepath     string
}

func NewLoader(validator *validator.Validate, isGlobal bool) *Reader {
	var (
		filePath      string
		fileExtension string
	)

	configFileName := os.Getenv(EnvironmentVariableConfigFileName)
	if isGlobal {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			zap.L().Fatal("failed to get current user home directory")
		}

		filePath = filepath.Join(homeDir, GlobalConfigFileName)
		fileExtension = ".json"
	} else {
		if configFileName == "" {
			configFileName = DefaultConfigFileName
		}

		workingDirectory, err := os.Getwd()
		if err != nil {
			zap.L().Fatal("failed to get current working directory path", zap.Error(err))
		}

		filePath = filepath.Join(workingDirectory, configFileName)
		fileExtension = filepath.Ext(configFileName)
		if fileExtension == "" {
			configFileName = filepath.Clean(fmt.Sprintf("%s.json", configFileName))
		}
	}

	configReader, err := NewReaderStrategy(fileExtension)
	if err != nil {
		zap.L().Fatal("failed to create config reader strategy", zap.Error(err))
	}

	return &Reader{
		ConfigReader: configReader,
		Validator:    validator,
		Filepath:     filePath,
	}
}

func (reader *Reader) Read() (*Cfg, error) {
	config, err := reader.ConfigReader.Read(reader.Filepath)
	if err != nil {
		return nil, err
	}

	if reader.Validator != nil {
		if err = reader.Validator.Struct(config); err != nil {
			return nil, err
		}
	}

	for routeIndex, route := range config.Routes {
		var (
			RequestToMethod = route.RequestTo.Method
			RequestToBody   = route.RequestTo.Body
		)

		if RequestToMethod == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if RequestToMethod == http.MethodGet && RequestToBody != nil {
			return nil, ErrorGetSendBody
		}

		config.Routes[routeIndex].Method = route.Method
		config.Routes[routeIndex].RequestTo.Method = route.RequestTo.Method
	}

	if config.HealthCheckRoute {
		config.Routes = append(
			config.Routes,
			Route{
				Method: "GET",
				Path:   "/health",
				FakeResponse: FakeResponse{
					StatusCode: http.StatusOK,
				},
			},
		)
	}

	if config.Concurrency == 0 {
		config.Concurrency = runtime.GOMAXPROCS(3)
	}

	return config, nil
}

func Write(config *Route, dir string) error {
	file, err := os.Create(filepath.Join(dir, "inzibat.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}

func InitGlobalConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get current user's home directory: %w", err)
	}

	globalConfigFilePath := filepath.Join(homeDir, DefaultConfigFileName)

	if _, err = os.Stat(globalConfigFilePath); !errors.Is(err, os.ErrNotExist) {
		var globalCfg *Cfg
		globalCfg, err = ReadOrCreateConfig(globalConfigFilePath)
		if err != nil {
			return fmt.Errorf("failed to read global config: %w", err)
		}

		if err = WriteConfig(globalCfg, globalConfigFilePath); err != nil {
			return fmt.Errorf("failed to initialize global config: %w", err)
		}
	}

	return nil
}

func ReadOrCreateConfig(configPath string) (*Cfg, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Cfg{
			ServerPort:       8080,
			Concurrency:      5,
			HealthCheckRoute: false,
			Routes:           []Route{},
		}, nil
	}

	ext := filepath.Ext(configPath)
	if ext == "" {
		ext = ".json"
	}

	readerStrategy, err := NewReaderStrategy(ext)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader strategy: %w", err)
	}

	cfg, err := readerStrategy.Read(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cfg, nil
}

func WriteConfig(cfg *Cfg, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}
