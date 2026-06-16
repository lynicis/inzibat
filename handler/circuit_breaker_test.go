package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lynicis/inzibat/config"
)

func TestBuildCircuitBreakerRouteKey(t *testing.T) {
	t.Run("happy path - without request to", func(t *testing.T) {
		route := config.Route{
			Method: http.MethodGet,
			Path:   "/proxy",
		}

		routeKey := BuildCircuitBreakerRouteKey(route)

		assert.Equal(t, "GET /proxy", routeKey)
	})

	t.Run("happy path - with request to", func(t *testing.T) {
		route := config.Route{
			Method: http.MethodGet,
			Path:   "/proxy",
			RequestTo: &config.RequestTo{
				Method: http.MethodPost,
				Host:   "http://example.com",
				Path:   "/api",
			},
		}

		routeKey := BuildCircuitBreakerRouteKey(route)

		assert.Equal(t, "GET /proxy -> POST http://example.com/api", routeKey)
	})
}

func TestCircuitBreakerStore_Seed(t *testing.T) {
	t.Run("happy path - update existing record", func(t *testing.T) {
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)

		routeKey := "GET /proxy"
		initialConfig := config.CircuitBreakerConfig{
			Enabled:          config.BoolPointer(true),
			FailureThreshold: 2,
		}
		updatedConfig := config.CircuitBreakerConfig{
			Enabled:          config.BoolPointer(true),
			FailureThreshold: 5,
		}

		assert.NoError(t, breakerStore.Seed(routeKey, initialConfig))
		assert.NoError(t, breakerStore.Seed(routeKey, updatedConfig))

		record, err := breakerStore.Get(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, 5, record.Config.FailureThreshold)
	})
}

func TestCircuitBreakerStore_Allow(t *testing.T) {
	t.Run("happy path - default state allows request", func(t *testing.T) {
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)

		routeKey := "GET /proxy"
		assert.NoError(t, breakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:             config.BoolPointer(true),
			FailureThreshold:    1,
			MinimumRequests:     1,
			OpenTimeoutMs:       10,
			HalfOpenMaxRequests: 1,
			SuccessThreshold:    1,
		}))

		record, err := breakerStore.Get(routeKey)
		assert.NoError(t, err)
		record.State = CircuitBreakerState("unknown")
		assert.NoError(t, breakerStore.update(routeKey, func(r *CircuitBreakerRecord) {
			*r = *record
		}))

		allowed, err := breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.True(t, allowed)
	})
}

func TestCircuitBreakerStore_OnSuccess(t *testing.T) {
	t.Run("happy path - closed state increments request count", func(t *testing.T) {
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)

		routeKey := "GET /proxy"
		assert.NoError(t, breakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:          config.BoolPointer(true),
			FailureThreshold: 2,
			MinimumRequests:  2,
		}))

		for range 3 {
			assert.NoError(t, breakerStore.OnSuccess(routeKey))
		}

		record, err := breakerStore.Get(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateClosed, record.State)
		assert.Equal(t, 3, record.RequestCount)
		assert.Equal(t, 0, record.ConsecutiveFailures)
	})
}

func TestCircuitBreakerStore_OnFailure(t *testing.T) {
	t.Run("happy path - open state does not change", func(t *testing.T) {
		now := time.Now()
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)
		breakerStore.clock = func() time.Time { return now }

		routeKey := "GET /proxy"
		assert.NoError(t, breakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:          config.BoolPointer(true),
			FailureThreshold: 1,
			MinimumRequests:  1,
		}))

		assert.NoError(t, breakerStore.OnFailure(routeKey))

		record, err := breakerStore.Get(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateOpen, record.State)

		assert.NoError(t, breakerStore.OnFailure(routeKey))

		record, err = breakerStore.Get(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateOpen, record.State)
	})

	t.Run("happy path - half-open state opens again", func(t *testing.T) {
		now := time.Now()
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)
		breakerStore.clock = func() time.Time { return now }

		routeKey := "GET /proxy"
		assert.NoError(t, breakerStore.Seed(routeKey, config.CircuitBreakerConfig{
			Enabled:             config.BoolPointer(true),
			FailureThreshold:    1,
			MinimumRequests:     1,
			OpenTimeoutMs:       10,
			HalfOpenMaxRequests: 1,
			SuccessThreshold:    1,
		}))

		assert.NoError(t, breakerStore.OnFailure(routeKey))
		now = now.Add(20 * time.Millisecond)
		allowed, err := breakerStore.Allow(routeKey)
		assert.NoError(t, err)
		assert.True(t, allowed)

		assert.NoError(t, breakerStore.OnFailure(routeKey))

		record, err := breakerStore.Get(routeKey)
		assert.NoError(t, err)
		assert.Equal(t, CircuitBreakerStateOpen, record.State)
	})
}

func TestCircuitBreakerStore_State(t *testing.T) {
	t.Run("error path - record not found", func(t *testing.T) {
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)

		state, err := breakerStore.State("non-existent")
		assert.Error(t, err)
		assert.Empty(t, state)
	})
}

func TestCircuitBreakerStore_Get(t *testing.T) {
	t.Run("error path - record not found", func(t *testing.T) {
		breakerStore, err := NewCircuitBreakerStore()
		assert.NoError(t, err)

		record, err := breakerStore.Get("non-existent")
		assert.Error(t, err)
		assert.Nil(t, record)
	})
}

func TestAllowHalfOpen(t *testing.T) {
	t.Run("happy path - non-half-open state returns false", func(t *testing.T) {
		record := &CircuitBreakerRecord{
			State: CircuitBreakerStateClosed,
		}

		allowed := allowHalfOpen(record)

		assert.False(t, allowed)
	})
}

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
