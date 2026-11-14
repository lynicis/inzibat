package config

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-json"
)

func readFile(filePath string) ([]byte, error) {
	absPath, err := ResolveAbsolutePath(filePath)
	if err != nil {
		return nil, err
	}

	// #nosec G304 - the path is resolved and cleaned by ResolveAbsolutePath
	file, err := os.Open(absPath)
	if err != nil {
		return nil, newFailOpeningError(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, newFailReadingError(err)
	}

	return data, nil
}

func LoadHeadersFromFile(filePath string) (http.Header, error) {
	data, err := readFile(filePath)
	if err != nil {
		return nil, err
	}

	var headersMap map[string]string
	if err := json.Unmarshal(data, &headersMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	headers := make(http.Header, len(headersMap))
	for key, value := range headersMap {
		headers.Set(key, value)
	}

	return headers, nil
}

func LoadBodyFromFile(filePath string) (HttpBody, error) {
	data, err := readFile(filePath)
	if err != nil {
		return nil, err
	}

	var body HttpBody
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return body, nil
}

func LoadBodyStringFromFile(filePath string) (string, error) {
	data, err := readFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
