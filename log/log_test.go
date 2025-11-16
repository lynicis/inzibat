package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLogInitialization(t *testing.T) {
	t.Run("happy path - logger is initialized", func(t *testing.T) {
		logger := zap.L()
		assert.NotNil(t, logger)
	})

	t.Run("happy path - logger can log messages", func(t *testing.T) {
		logger := zap.L()

		assert.NotPanics(t, func() {
			logger.Info("test message")
		})
	})

	t.Run("happy path - logger has structured logging", func(t *testing.T) {
		logger := zap.L()

		assert.NotPanics(t, func() {
			logger.Info("test message",
				zap.String("key", "value"),
				zap.Int("number", 42),
			)
		})
	})
}
