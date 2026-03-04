package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-memdb"

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
	HalfOpenRequests     int
	UpdatedAt            time.Time
}

type CircuitBreakerStore struct {
	db    *memdb.MemDB
	clock func() time.Time
	mu    sync.Mutex
}

func NewCircuitBreakerStore() (*CircuitBreakerStore, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"circuit_breakers": {
				Name: "circuit_breakers",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "RouteKey"},
					},
					"state": {
						Name:    "state",
						Indexer: &memdb.StringFieldIndex{Field: "State"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit breaker memdb: %w", err)
	}

	return &CircuitBreakerStore{db: db, clock: time.Now}, nil
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

	txn := store.db.Txn(true)
	defer txn.Abort()

	existingRaw, err := txn.First("circuit_breakers", "id", routeKey)
	if err != nil {
		return fmt.Errorf("failed to lookup circuit breaker record: %w", err)
	}

	if existingRaw != nil {
		existingRecord := existingRaw.(*CircuitBreakerRecord)
		existingRecord.Config = cfg
		existingRecord.UpdatedAt = store.clock()
		if err = txn.Insert("circuit_breakers", existingRecord); err != nil {
			return fmt.Errorf("failed to update circuit breaker record: %w", err)
		}
		txn.Commit()
		return nil
	}

	record := &CircuitBreakerRecord{
		RouteKey:  routeKey,
		State:     CircuitBreakerStateClosed,
		Config:    cfg,
		UpdatedAt: store.clock(),
	}
	if err = txn.Insert("circuit_breakers", record); err != nil {
		return fmt.Errorf("failed to insert circuit breaker record: %w", err)
	}

	txn.Commit()
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
				record.HalfOpenRequests = 0
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
	txn := store.db.Txn(false)
	defer txn.Abort()

	recordRaw, err := txn.First("circuit_breakers", "id", routeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup circuit breaker record: %w", err)
	}

	if recordRaw == nil {
		return nil, fmt.Errorf("circuit breaker record not found for route key %s", routeKey)
	}

	record := recordRaw.(*CircuitBreakerRecord)
	recordCopy := *record
	return &recordCopy, nil
}

func (store *CircuitBreakerStore) update(routeKey string, updateFn func(record *CircuitBreakerRecord)) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	txn := store.db.Txn(true)
	defer txn.Abort()

	recordRaw, err := txn.First("circuit_breakers", "id", routeKey)
	if err != nil {
		return fmt.Errorf("failed to lookup circuit breaker record: %w", err)
	}

	if recordRaw == nil {
		return fmt.Errorf("circuit breaker record not found for route key %s", routeKey)
	}

	record := recordRaw.(*CircuitBreakerRecord)
	updateFn(record)
	record.UpdatedAt = store.clock()

	if err = txn.Insert("circuit_breakers", record); err != nil {
		return fmt.Errorf("failed to update circuit breaker record: %w", err)
	}

	txn.Commit()
	return nil
}

func allowHalfOpen(record *CircuitBreakerRecord) bool {
	if record.State != CircuitBreakerStateHalfOpen {
		return false
	}

	if record.HalfOpenRequests >= record.Config.HalfOpenMaxRequests {
		return false
	}

	record.HalfOpenRequests++
	return true
}

func openRecord(record *CircuitBreakerRecord, now time.Time) {
	record.State = CircuitBreakerStateOpen
	record.OpenedAt = now
	record.ConsecutiveFailures = 0
	record.ConsecutiveSuccesses = 0
	record.HalfOpenRequests = 0
}

func resetToClosed(record *CircuitBreakerRecord) {
	record.State = CircuitBreakerStateClosed
	record.RequestCount = 0
	record.ConsecutiveFailures = 0
	record.ConsecutiveSuccesses = 0
	record.HalfOpenRequests = 0
}
