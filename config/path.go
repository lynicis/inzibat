package config

import (
	"fmt"
	"path/filepath"
)

func ResolveAbsolutePath(filePath string) (string, error) {
	cleanPath := filepath.Clean(filePath)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}
	return absPath, nil
}

