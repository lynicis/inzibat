package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveAbsolutePath(t *testing.T) {
	t.Run("happy path - resolves relative path", func(t *testing.T) {
		relativePath := "./test.json"
		expectedAbs, err := filepath.Abs(relativePath)
		require.NoError(t, err)

		absPath, err := ResolveAbsolutePath(relativePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedAbs, absPath)
	})

	t.Run("happy path - resolves absolute path", func(t *testing.T) {
		tmpDir := t.TempDir()
		absPath := filepath.Join(tmpDir, "test.json")

		resolved, err := ResolveAbsolutePath(absPath)

		assert.NoError(t, err)
		assert.Equal(t, absPath, resolved)
	})

	t.Run("happy path - cleans path with ..", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.json")
		err := os.WriteFile(testFile, []byte("{}"), 0644)
		require.NoError(t, err)

		relativePath := filepath.Join(tmpDir, "..", filepath.Base(tmpDir), "test.json")
		expectedAbs := filepath.Join(tmpDir, "test.json")

		resolved, err := ResolveAbsolutePath(relativePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedAbs, resolved)
	})

	t.Run("happy path - handles current directory", func(t *testing.T) {
		wd, err := os.Getwd()
		require.NoError(t, err)
		relativePath := "test.json"
		expectedAbs := filepath.Join(wd, "test.json")

		resolved, err := ResolveAbsolutePath(relativePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedAbs, resolved)
	})

	t.Run("error path - invalid path that cannot be resolved", func(t *testing.T) {
		veryLongPath := string(make([]byte, 4096)) + "/test.json"

		_, err := ResolveAbsolutePath(veryLongPath)

		if err != nil {
			assert.Contains(t, err.Error(), "failed to resolve file path")
		}
	})
}
