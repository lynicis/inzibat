package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/lynicis/inzibat/config"
)

type CircuitBreakerState string

const (
	CircuitBreakerStateClosed   CircuitBreakerState = "closed"
	CircuitBreakerStateOpen     CircuitBreakerState = "open"
	CircuitBreakerStateHalfOpen CircuitBreakerState = "half-open"
)

type CircuitBreakerRecord struct {
	RouteKey             string
	State                CircuitBreakerState
	Config               config.CircuitBreakerConfig
	OpenedAt             time.Time
	RequestCount         int
	ConsecutiveFailures  int
	ConsecutiveSuccesses int
	HalfOpenInFlight     int
	UpdatedAt            time.Time
}

type CircuitBreakerStore struct {
	records map[string]*CircuitBreakerRecord
	clock   func() time.Time
	mu      sync.Mutex
}

func NewCircuitBreakerStore() *CircuitBreakerStore {
	return &CircuitBreakerStore{
		records: make(map[string]*CircuitBreakerRecord),
		clock:   time.Now,
	}
}

func BuildCircuitBreakerRouteKey(route config.Route) string {
	if route.RequestTo == nil {
		return fmt.Sprintf("%s %s", route.Method, route.Path)
	}

	return fmt.Sprintf(
		"%s %s -> %s %s%s",
		route.Method,
		route.Path,
		route.RequestTo.Method,
		route.RequestTo.Host,
		route.RequestTo.Path,
	)
}

func (store *CircuitBreakerStore) Seed(routeKey string, cfg config.CircuitBreakerConfig) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if record, ok := store.records[routeKey]; ok {
		record.Config = cfg
		record.UpdatedAt = store.clock()

		return nil
	}

	store.records[routeKey] = &CircuitBreakerRecord{
		RouteKey:  routeKey,
		State:     CircuitBreakerStateClosed,
		Config:    cfg,
		UpdatedAt: store.clock(),
	}

	return nil
}

func (store *CircuitBreakerStore) Allow(routeKey string) (bool, error) {
	allowed := false
	err := store.update(routeKey, func(record *CircuitBreakerRecord) {
		switch record.State {
		case CircuitBreakerStateClosed:
			allowed = true
		case CircuitBreakerStateOpen:
			openTimeout := time.Duration(record.Config.OpenTimeoutMs) * time.Millisecond
			if store.clock().Sub(record.OpenedAt) >= openTimeout {
				record.State = CircuitBreakerStateHalfOpen
				record.HalfOpenInFlight = 0
				record.ConsecutiveSuccesses = 0
			}
			allowed = allowHalfOpen(record)
		case CircuitBreakerStateHalfOpen:
			allowed = allowHalfOpen(record)
		default:
			allowed = true
		}
	})

	return allowed, err
}

func (store *CircuitBreakerStore) OnSuccess(routeKey string) error {
	return store.update(routeKey, func(record *CircuitBreakerRecord) {
		if record.State == CircuitBreakerStateHalfOpen {
			releaseHalfOpen(record)
			record.ConsecutiveSuccesses++
			if record.ConsecutiveSuccesses >= record.Config.SuccessThreshold {
				resetToClosed(record)
			}
			return
		}

		if record.State == CircuitBreakerStateClosed {
			record.RequestCount++
			record.ConsecutiveFailures = 0
		}
	})
}

func (store *CircuitBreakerStore) OnFailure(routeKey string) error {
	return store.update(routeKey, func(record *CircuitBreakerRecord) {
		if record.State == CircuitBreakerStateHalfOpen {
			releaseHalfOpen(record)
			openRecord(record, store.clock())
			return
		}

		if record.State == CircuitBreakerStateOpen {
			return
		}

		record.RequestCount++
		record.ConsecutiveFailures++
		if record.RequestCount >= record.Config.MinimumRequests &&
			record.ConsecutiveFailures >= record.Config.FailureThreshold {
			openRecord(record, store.clock())
		}
	})
}

func (store *CircuitBreakerStore) State(routeKey string) (CircuitBreakerState, error) {
	record, err := store.Get(routeKey)
	if err != nil {
		return "", err
	}

	return record.State, nil
}

func (store *CircuitBreakerStore) Get(routeKey string) (*CircuitBreakerRecord, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	record, ok := store.records[routeKey]
	if !ok {
		return nil, fmt.Errorf("circuit breaker record not found for route key %s", routeKey)
	}

	recordCopy := *record

	return &recordCopy, nil
}

func (store *CircuitBreakerStore) update(routeKey string, updateFn func(record *CircuitBreakerRecord)) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	record, ok := store.records[routeKey]
	if !ok {
		return fmt.Errorf("circuit breaker record not found for route key %s", routeKey)
	}

	updateFn(record)
	record.UpdatedAt = store.clock()

	return nil
}

// releaseHalfOpen frees a half-open slot held by the request that just finished.
func releaseHalfOpen(record *CircuitBreakerRecord) {
	if record.HalfOpenInFlight > 0 {
		record.HalfOpenInFlight--
	}
}

func allowHalfOpen(record *CircuitBreakerRecord) bool {
	if record.State != CircuitBreakerStateHalfOpen {
		return false
	}

	if record.HalfOpenInFlight >= record.Config.HalfOpenMaxRequests {
		return false
	}

	record.HalfOpenInFlight++
	return true
}

func openRecord(record *CircuitBreakerRecord, now time.Time) {
	record.State = CircuitBreakerStateOpen
	record.OpenedAt = now
	record.ConsecutiveFailures = 0
	record.ConsecutiveSuccesses = 0
	record.HalfOpenInFlight = 0
}

func resetToClosed(record *CircuitBreakerRecord) {
	record.State = CircuitBreakerStateClosed
	record.RequestCount = 0
	record.ConsecutiveFailures = 0
	record.ConsecutiveSuccesses = 0
	record.HalfOpenInFlight = 0
}
