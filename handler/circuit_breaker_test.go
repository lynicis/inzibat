package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lynicis/inzibat/config"
)

func TestCircuitBreaker(t *testing.T) {
	t.Run("opens after threshold is reached", func(t *testing.T) {
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)

		routeKey := "GET /proxy -> GET http://example.com/api"
		err = breakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:             config.BoolPointer(true),
			FailureThreshold:    2,
			MinimumRequests:     2,
			OpenTimeoutMs:       1000,
			HalfOpenMaxRequests: 1,
			SuccessThreshold:    1,
		})
		assert.NoError(t, err)

		allowed, err := breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.True(t, allowed)
		assert.NoError(t, breakerStore.OnFailure(routeKey))

		allowed, err = breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.True(t, allowed)
		assert.NoError(t, breakerStore.OnFailure(routeKey))

		state, err := breakerStore.State(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateOpen, state)

		allowed, err = breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("transitions to half-open and closes after success threshold", func(t *testing.T) {
		now := time.Now()
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)
		breakerStore.clock = func() time.Time { return now }

		routeKey := "GET /proxy -> GET http://example.com/api"
		err = breakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:             config.BoolPointer(true),
			FailureThreshold:    1,
			MinimumRequests:     1,
			OpenTimeoutMs:       10,
			HalfOpenMaxRequests: 2,
			SuccessThreshold:    2,
		})
		assert.NoError(t, err)

		assert.NoError(t, breakerStore.OnFailure(routeKey))
		state, err := breakerStore.State(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateOpen, state)

		now = now.Add(20 * time.Millisecond)
		allowed, err := breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.True(t, allowed)

		state, err = breakerStore.State(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateHalfOpen, state)
		assert.NoError(t, breakerStore.OnSuccess(routeKey))

		allowed, err = breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.True(t, allowed)
		assert.NoError(t, breakerStore.OnSuccess(routeKey))

		state, err = breakerStore.State(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateClosed, state)
	})
}
