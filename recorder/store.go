package recorder

import (
	"sync"
	"time"
)

// Store is a thread-safe, capacity-limited, in-memory store for recorded entries.
// When the capacity is exceeded, oldest entries are dropped (ring-buffer behavior).
type Store struct {
	mu        sync.RWMutex
	entries   []RecordedEntry
	capacity  int
	startedAt time.Time
}

// NewStore creates a new Store with the given capacity.
func NewStore(capacity int) *Store {
	if capacity <= 0 {
		capacity = DefaultStoreCapacity
	}

	return &Store{
		entries:   make([]RecordedEntry, 0, min(capacity, 1024)),
		capacity:  capacity,
		startedAt: time.Now(),
	}
}

// Add appends a recorded entry. If the store is at capacity, the oldest entry is dropped.
func (s *Store) Add(entry RecordedEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.entries) >= s.capacity {
		s.entries = s.entries[1:]
	}

	s.entries = append(s.entries, entry)
}

// List returns a copy of all recorded entries.
func (s *Store) List() []RecordedEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]RecordedEntry, len(s.entries))
	copy(result, s.entries)

	return result
}

// Session returns the full recorded session with metadata.
func (s *Store) Session() RecordedSession {
	entries := s.List()

	return RecordedSession{
		StartedAt:  s.startedAt,
		EntryCount: len(entries),
		Entries:    entries,
	}
}

// Clear removes all recorded entries.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = make([]RecordedEntry, 0, min(s.capacity, 1024))
}

// Len returns the number of recorded entries.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.entries)
}
