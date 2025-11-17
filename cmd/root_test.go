package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("happy path - root command has correct properties", func(t *testing.T) {
		assert.NotNil(t, rootCmd)
		assert.Equal(t, "inzibat", rootCmd.Use)
		assert.Contains(t, rootCmd.Short, "HTTP mock server")
		assert.Contains(t, rootCmd.Long, "Inzibat")
		assert.Contains(t, rootCmd.Long, "Military Police")
		assert.Contains(t, rootCmd.Long, "JSON/TOML/YAML")
	})

	t.Run("happy path - Execute() completes successfully with help flag", func(t *testing.T) {
		originalArgs := os.Args
		defer func() {
			os.Args = originalArgs
		}()

		os.Args = []string{"inzibat", "--help"}

		err := rootCmd.Execute()

		assert.NoError(t, err, "Execute() should complete successfully with --help flag")
	})

	t.Run("happy path - Execute() function completes without exit on success", func(t *testing.T) {
		originalArgs := os.Args
		defer func() {
			os.Args = originalArgs
		}()

		os.Args = []string{"inzibat", "--help"}

		Execute()

		assert.True(t, true, "Execute() completed successfully without exit")
	})

	t.Run("happy path - rootCmd has all expected subcommands", func(t *testing.T) {
		assert.NotNil(t, rootCmd)
		commands := rootCmd.Commands()
		commandNames := make(map[string]bool)
		for _, cmd := range commands {
			commandNames[cmd.Use] = true
		}
		assert.True(t, commandNames["list"], "list command should be registered")
		assert.True(t, commandNames["start"], "start command should be registered")
		assert.True(t, commandNames["create"], "create command should be registered")
	})
}

func TestExecute_ErrorPath(t *testing.T) {
	t.Run("error path - exits with code 1 on invalid command", func(t *testing.T) {
		_, testFile, _, ok := runtime.Caller(0)
		require.True(t, ok, "failed to get test file path")

		testDir := filepath.Dir(testFile)
		projectRoot := filepath.Dir(testDir)
		projectRoot, err := filepath.Abs(projectRoot)
		require.NoError(t, err)

		binaryPath := filepath.Join(projectRoot, "inzibat")

		if _, statErr := os.Stat(binaryPath); os.IsNotExist(statErr) {
			buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
			buildCmd.Dir = projectRoot
			buildErr := buildCmd.Run()
			require.NoError(t, buildErr, "failed to build binary")
		}

		info, statErr := os.Stat(binaryPath)
		if statErr != nil {
			t.Fatalf("binary does not exist at %s: %v", binaryPath, statErr)
		}
		if info.IsDir() {
			t.Fatalf("binary path %s is a directory, not a file. projectRoot: %s", binaryPath, projectRoot)
		}
		if info.Mode()&0111 == 0 {
			t.Fatalf("binary at %s is not executable", binaryPath)
		}

		cmd := exec.Command(binaryPath, "invalid-command-that-does-not-exist")
		err = cmd.Run()

		if exitError, ok := err.(*exec.ExitError); ok {
			assert.Equal(t, 1, exitError.ExitCode(), "Execute() should exit with code 1 on error")
		} else if err != nil {
			t.Fatalf("unexpected error type: %v (expected ExitError)", err)
		} else {
			t.Fatalf("expected command to fail with exit code 1, but it succeeded")
		}
	})

	t.Run("error path - Execute() calls exitFunc with code 1 when rootCmd.Execute() returns error", func(t *testing.T) {
		originalExitFunc := exitFunc
		originalRootCmd := rootCmd
		defer func() {
			exitFunc = originalExitFunc
			rootCmd = originalRootCmd
		}()

		var capturedExitCode int
		exitFunc = func(code int) {
			capturedExitCode = code
		}

		errorCmd := &cobra.Command{
			Use: "test-error-cmd",
			RunE: func(cmd *cobra.Command, args []string) error {
				return errors.New("test error")
			},
		}
		rootCmd = errorCmd

		Execute()

		assert.Equal(t, 1, capturedExitCode, "Execute() should call exitFunc with code 1 when rootCmd.Execute() returns error")
	})
}
