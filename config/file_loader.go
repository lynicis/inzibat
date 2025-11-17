package config

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-json"
)

type FileLoader interface {
	Load(filePath string) (interface{}, error)
}

type HeadersLoader struct{}

func (l *HeadersLoader) Load(filePath string) (interface{}, error) {
	absPath, err := ResolveAbsolutePath(filePath)
	if err != nil {
		return nil, err
	}
	// #nosec G304 - File path is validated and cleaned before use
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var headersMap map[string]string
	if err := json.Unmarshal(data, &headersMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	headers := make(http.Header)
	for key, value := range headersMap {
		headers.Set(key, value)
	}

	return headers, nil
}

type BodyLoader struct{}

func (l *BodyLoader) Load(filePath string) (interface{}, error) {
	absPath, err := ResolveAbsolutePath(filePath)
	if err != nil {
		return nil, err
	}
	// #nosec G304 - File path is validated and cleaned before use
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var body HttpBody
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return body, nil
}

type BodyStringLoader struct{}

func (l *BodyStringLoader) Load(filePath string) (interface{}, error) {
	absPath, err := ResolveAbsolutePath(filePath)
	if err != nil {
		return "", err
	}
	// #nosec G304
	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}

func loadFromFile(filePath string, loader FileLoader) (interface{}, error) {
	return loader.Load(filePath)
}

func LoadHeadersFromFile(filePath string) (http.Header, error) {
	loader := &HeadersLoader{}
	result, err := loadFromFile(filePath, loader)
	if err != nil {
		return nil, err
	}
	return result.(http.Header), nil
}

func LoadBodyFromFile(filePath string) (HttpBody, error) {
	loader := &BodyLoader{}
	result, err := loadFromFile(filePath, loader)
	if err != nil {
		return nil, err
	}
	return result.(HttpBody), nil
}

func LoadBodyStringFromFile(filePath string) (string, error) {
	loader := &BodyStringLoader{}
	result, err := loadFromFile(filePath, loader)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
