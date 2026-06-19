package recorder

import (
	"time"

	"github.com/goccy/go-json"
)

// MaxBodyCaptureBytes is the maximum size of request/response body that will be captured.
const MaxBodyCaptureBytes = 1 << 20 // 1 MB

// DefaultStoreCapacity is the maximum number of entries the store will hold.
const DefaultStoreCapacity = 10_000

// RecordedEntry represents a single captured request-response pair.
type RecordedEntry struct {
	ID         string           `json:"id"`
	Timestamp  time.Time        `json:"timestamp"`
	Request    RecordedRequest  `json:"request"`
	Response   RecordedResponse `json:"response"`
	DurationMs int64            `json:"durationMs"`
}

// RecordedRequest captures the incoming HTTP request metadata.
type RecordedRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    json.RawMessage     `json:"body,omitempty"`
}

// RecordedResponse captures the outgoing HTTP response metadata.
type RecordedResponse struct {
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Body       json.RawMessage     `json:"body,omitempty"`
}

// RecordedSession wraps a set of recorded entries with session metadata.
type RecordedSession struct {
	StartedAt  time.Time       `json:"startedAt"`
	EntryCount int             `json:"entryCount"`
	Entries    []RecordedEntry `json:"entries"`
}
