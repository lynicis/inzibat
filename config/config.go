package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"

	"github.com/go-playground/validator/v10"
)

const DefaultConfigExtension = ".json"

type Reader struct {
	ConfigReader ReaderStrategy
	Validator    *validator.Validate
	Filepath     string
}

func resolveConfigPath(isGlobal bool, explicitPath string) (string, string, error) {
	if explicitPath != "" {
		extension := filepath.Ext(explicitPath)
		if extension == "" {
			extension = DefaultConfigExtension
		}

		return explicitPath, extension, nil
	}

	if isGlobal {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", "", fmt.Errorf("failed to get current user home directory: %w", err)
		}

		return filepath.Join(homeDir, GlobalConfigFileName), DefaultConfigExtension, nil
	}

	configFileName := os.Getenv(EnvironmentVariableConfigFileName)
	if configFileName == "" {
		configFileName = DefaultConfigFileName
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current working directory path: %w", err)
	}

	extension := filepath.Ext(configFileName)
	if extension == "" {
		configFileName = filepath.Clean(fmt.Sprintf("%s.json", configFileName))
		extension = DefaultConfigExtension
	}

	return filepath.Join(workingDirectory, configFileName), extension, nil
}

// ResolveConfigPath returns the config file path the commands would use.
func ResolveConfigPath(isGlobal bool, explicitPath string) (string, error) {
	filePath, _, err := resolveConfigPath(isGlobal, explicitPath)
	return filePath, err
}

func NewLoader(validator *validator.Validate, isGlobal bool, explicitPath string) (*Reader, error) {
	filePath, fileExtension, err := resolveConfigPath(isGlobal, explicitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config file path: %w", err)
	}

	configReader, err := NewReaderStrategy(fileExtension)
	if err != nil {
		return nil, fmt.Errorf("failed to create config reader strategy: %w", err)
	}

	return &Reader{
		ConfigReader: configReader,
		Validator:    validator,
		Filepath:     filePath,
	}, nil
}

func (reader *Reader) Read() (*Cfg, error) {
	config, err := reader.ConfigReader.Read(reader.Filepath)
	if err != nil {
		return nil, err
	}

	if err = reader.validate(config); err != nil {
		return nil, err
	}

	if config.CircuitBreaker != nil {
		config.CircuitBreaker = MergeCircuitBreakerConfig(nil, config.CircuitBreaker)
	}

	if err = normalizeRoutes(config); err != nil {
		return nil, err
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

	return config, nil
}

func (reader *Reader) validate(config *Cfg) error {
	if reader.Validator == nil {
		return nil
	}

	return reader.Validator.Struct(config)
}

func normalizeRoutes(config *Cfg) error {
	for routeIndex := range config.Routes {
		route := &config.Routes[routeIndex]
		if route.RequestTo == nil {
			continue
		}

		route.RequestTo.CircuitBreaker = MergeCircuitBreakerConfig(
			config.CircuitBreaker,
			route.RequestTo.CircuitBreaker,
		)

		if route.RequestTo.Method == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if route.RequestTo.Method == http.MethodGet && route.RequestTo.Body != nil {
			return ErrorGetSendBody
		}
	}

	return nil
}

func WriteConfig(cfg *Cfg, filePath string) error {
	if cfg.Schema == "" {
		cfg.Schema = DefaultSchemaFileName
	}

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

	if err = json.NewEncoder(file).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if err = file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

func DefaultCfg() *Cfg {
	return &Cfg{
		Schema:           DefaultSchemaFileName,
		ServerPort:       8080,
		HealthCheckRoute: false,
		Routes:           []Route{},
	}
}

func ReadOrCreateConfig(configPath string) (*Cfg, error) {
	_, err := os.Stat(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		return DefaultCfg(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	_, extension, err := resolveConfigPath(false, configPath)
	if err != nil {
		return nil, err
	}

	readerStrategy, err := NewReaderStrategy(extension)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader strategy: %w", err)
	}

	cfg, err := readerStrategy.Read(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if cfg.Schema == "" {
		cfg.Schema = DefaultSchemaFileName
	}

	return cfg, nil
}
