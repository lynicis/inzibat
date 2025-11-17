package config

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goccy/go-json"

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

	for routeIndex := range config.Routes {
		route := &config.Routes[routeIndex]
		if route.RequestTo == nil {
			continue
		}

		if route.RequestTo.Method == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if route.RequestTo.Method == http.MethodGet && route.RequestTo.Body != nil {
			return nil, ErrorGetSendBody
		}
	}

	if config.HealthCheckRoute {
		config.Routes = append(
			config.Routes,
			Route{
				Method: "GET",
				Path:   "/health",
				FakeResponse: &FakeResponse{
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

func writeJSONToFile(filePath string, data interface{}) error {
	absPath, err := ResolveAbsolutePath(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve file path: %w", err)
	}
	// #nosec G304 - File path is validated and cleaned before use
	file, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if err = file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

func Write(config *Route, dir string) error {
	absDir, err := ResolveAbsolutePath(dir)
	if err != nil {
		return fmt.Errorf("failed to resolve directory path: %w", err)
	}
	return writeJSONToFile(filepath.Join(absDir, "inzibat.json"), config)
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

func WriteConfig(cfg *Cfg, filePath string) error {
	return writeJSONToFile(filePath, cfg)
}
