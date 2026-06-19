package recorder

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEntry(id string) RecordedEntry {
	return RecordedEntry{
		ID:        id,
		Timestamp: time.Now(),
		Request: RecordedRequest{
			Method: "GET",
			Path:   "/test/" + id,
		},
		Response: RecordedResponse{
			StatusCode: 200,
		},
		DurationMs: 10,
	}
}

func TestNewStore(t *testing.T) {
	t.Run("creates store with given capacity", func(t *testing.T) {
		store := NewStore(500)
		assert.Equal(t, 0, store.Len())
		assert.Equal(t, 500, store.capacity)
	})

	t.Run("uses default capacity for zero", func(t *testing.T) {
		store := NewStore(0)
		assert.Equal(t, DefaultStoreCapacity, store.capacity)
	})

	t.Run("uses default capacity for negative", func(t *testing.T) {
		store := NewStore(-1)
		assert.Equal(t, DefaultStoreCapacity, store.capacity)
	})
}

func TestStoreAdd(t *testing.T) {
	t.Run("adds entries", func(t *testing.T) {
		store := NewStore(10)
		store.Add(newTestEntry("1"))
		store.Add(newTestEntry("2"))
		assert.Equal(t, 2, store.Len())
	})

	t.Run("drops oldest when at capacity", func(t *testing.T) {
		store := NewStore(3)
		store.Add(newTestEntry("1"))
		store.Add(newTestEntry("2"))
		store.Add(newTestEntry("3"))
		store.Add(newTestEntry("4"))

		assert.Equal(t, 3, store.Len())

		entries := store.List()
		assert.Equal(t, "2", entries[0].ID)
		assert.Equal(t, "3", entries[1].ID)
		assert.Equal(t, "4", entries[2].ID)
	})
}

func TestStoreList(t *testing.T) {
	t.Run("returns copy of entries", func(t *testing.T) {
		store := NewStore(10)
		store.Add(newTestEntry("1"))

		entries := store.List()
		require.Len(t, entries, 1)
		assert.Equal(t, "1", entries[0].ID)

		// Modify returned slice — store should be unaffected
		entries[0].ID = "modified"
		originalEntries := store.List()
		assert.Equal(t, "1", originalEntries[0].ID)
	})

	t.Run("returns empty slice when empty", func(t *testing.T) {
		store := NewStore(10)
		entries := store.List()
		assert.Empty(t, entries)
		assert.NotNil(t, entries)
	})
}

func TestStoreSession(t *testing.T) {
	store := NewStore(10)
	store.Add(newTestEntry("1"))
	store.Add(newTestEntry("2"))

	session := store.Session()
	assert.Equal(t, 2, session.EntryCount)
	assert.Len(t, session.Entries, 2)
	assert.False(t, session.StartedAt.IsZero())
}

func TestStoreClear(t *testing.T) {
	store := NewStore(10)
	store.Add(newTestEntry("1"))
	store.Add(newTestEntry("2"))
	assert.Equal(t, 2, store.Len())

	store.Clear()
	assert.Equal(t, 0, store.Len())
	assert.Empty(t, store.List())
}

func TestStoreConcurrency(t *testing.T) {
	store := NewStore(1000)

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Add(newTestEntry(fmt.Sprintf("%d", i)))
		}()
	}

	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.List()
		}()
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.Len()
		}()
	}

	wg.Wait()
	assert.LessOrEqual(t, store.Len(), 1000)
}
