package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lynicis/inzibat/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchRecordedEntries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/_inzibat/recorder/entries", r.URL.Path)
		entries := []recorder.RecordedEntry{
			{DurationMs: 100},
		}
		json.NewEncoder(w).Encode(entries)
	}))
	defer ts.Close()

	addr := ts.Listener.Addr().String()
	entries, err := fetchRecordedEntries(addr)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, int64(100), entries[0].DurationMs)

	// Test error cases
	_, err = fetchRecordedEntries("invalid-addr")
	require.Error(t, err)

	tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer tsErr.Close()
	_, err = fetchRecordedEntries(tsErr.Listener.Addr().String())
	require.Error(t, err)
}

func TestFetchRecordedSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/_inzibat/recorder/session", r.URL.Path)
		session := recorder.RecordedSession{
			StartedAt:  time.Now(),
		}
		json.NewEncoder(w).Encode(session)
	}))
	defer ts.Close()

	addr := ts.Listener.Addr().String()
	session, err := fetchRecordedSession(addr)
	require.NoError(t, err)
	require.NotNil(t, session)

	// Test error cases
	_, err = fetchRecordedSession("invalid-addr")
	require.Error(t, err)

	tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer tsErr.Close()
	_, err = fetchRecordedSession(tsErr.Listener.Addr().String())
	require.Error(t, err)
}

func TestExportAsJSON(t *testing.T) {
	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "export.json")

	oldOutput := recordExportOutput
	recordExportOutput = exportPath
	defer func() { recordExportOutput = oldOutput }()

	session := &recorder.RecordedSession{
		Entries: []recorder.RecordedEntry{},
	}

	exportAsJSON(session)

	data, err := os.ReadFile(exportPath)
	require.NoError(t, err)
	require.Contains(t, string(data), "startedAt")
}

func TestExportAsInzibatConfig(t *testing.T) {
	tempDir := t.TempDir()
	exportPath := filepath.Join(tempDir, "inzibat.json")

	oldOutput := recordExportOutput
	recordExportOutput = exportPath
	defer func() { recordExportOutput = oldOutput }()

	session := &recorder.RecordedSession{
		Entries: []recorder.RecordedEntry{},
	}

	exportAsInzibatConfig(session)

	data, err := os.ReadFile(exportPath)
	require.NoError(t, err)
	require.Contains(t, string(data), "routes")
}
