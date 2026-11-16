package formbuilder

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lynicis/inzibat/config"
)

// Key sequences for driving huh forms through a swapped stdin. Selects only
// accept arrow keys here: typing into a select's filter panics on a non-TTY
// output. Confirm fields advance on the bare y/n key, so no enter follows.
const (
	headlessEnter = "\r"
	headlessDown  = "\x1b[B"
	headlessYes   = "y"
	headlessNo    = "n"
)

// headlessSelect picks the zero-based index-th option of a Select field.
func headlessSelect(index int) string {
	return strings.Repeat(headlessDown, index) + headlessEnter
}

// withHeadlessTerminal runs fn with os.Stdin replaced by a pipe fed the key
// segments and os.Stdout sent to /dev/null, so real huh forms run headless.
// TERM is forced to a real value: huh silently switches to accessible mode
// under TERM=dumb, which never populates form results. Each segment is one
// field's keys, written separately: a form drops the enter between two
// fields when a whole multi-field stream arrives in a single read, and a
// previous form's reader slurps segments written before it exits. A regular
// file cannot stand in for the pipe: kqueue never signals EOF on it and the
// form blocks forever. Tests using it must not run in parallel.
func withHeadlessTerminal(t *testing.T, fn func(), segments ...string) {
	t.Helper()

	t.Setenv("TERM", "xterm-256color")

	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	require.NoError(t, err)
	t.Cleanup(func() { devnull.Close() })

	r, w, err := os.Pipe()
	require.NoError(t, err)

	oldIn, oldOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = r, devnull
	defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()
	defer func() { _ = r.Close() }()

	go func() {
		for _, segment := range segments {
			_, _ = w.WriteString(segment)
			time.Sleep(200 * time.Millisecond)
		}
		time.Sleep(200 * time.Millisecond)
		_ = w.Close()
	}()

	fn()
}

func TestCollectHeadersFromFormHeadless(t *testing.T) {
	var headers http.Header

	withHeadlessTerminal(t, func() {
		var err error
		headers, err = CollectHeadersFromForm()
		require.NoError(t, err)
	}, "X-Test\r", "hello\r", headlessNo)

	assert.Equal(t, "hello", headers.Get("X-Test"))
}

func TestCollectHeadersFromFormHeadlessMultiplePairs(t *testing.T) {
	var headers http.Header

	withHeadlessTerminal(t, func() {
		var err error
		headers, err = CollectHeadersFromForm()
		require.NoError(t, err)
	}, "X-A\r", "1\r", headlessYes, "X-B\r", "2\r", headlessNo)

	assert.Equal(t, "1", headers.Get("X-A"))
	assert.Equal(t, "2", headers.Get("X-B"))
	assert.Len(t, headers, 2)
}

func TestCollectBodyFromFormHeadless(t *testing.T) {
	var body config.HttpBody

	withHeadlessTerminal(t, func() {
		var err error
		body, err = CollectBodyFromForm()
		require.NoError(t, err)
	}, "message\r", "success\r", headlessNo)

	assert.Equal(t, config.HttpBody{"message": "success"}, body)
}

func TestCollectBodyStringFromFormHeadless(t *testing.T) {
	var bodyString string

	withHeadlessTerminal(t, func() {
		var err error
		bodyString, err = CollectBodyStringFromForm()
		require.NoError(t, err)
	}, "raw payload\r")

	assert.Equal(t, "raw payload", bodyString)
}

func TestCollectHeadersHeadlessFormSource(t *testing.T) {
	var headers http.Header

	withHeadlessTerminal(t, func() {
		var err error
		headers, err = CollectHeaders()
		require.NoError(t, err)
	}, headlessSelect(1), "X-Key\r", "from-form\r", headlessNo)

	assert.Equal(t, "from-form", headers.Get("X-Key"))
}

func TestCollectHeadersHeadlessFileSource(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "headers.json")
	require.NoError(t, os.WriteFile(filePath, []byte(`{"X-File": "loaded"}`), 0o600))

	var headers http.Header

	withHeadlessTerminal(t, func() {
		var err error
		headers, err = CollectHeaders()
		require.NoError(t, err)
	}, headlessSelect(0), filePath+headlessEnter)

	assert.Equal(t, "loaded", headers.Get("X-File"))
}

func TestCollectBodyHeadlessSkipSource(t *testing.T) {
	var body config.HttpBody

	withHeadlessTerminal(t, func() {
		var err error
		body, err = CollectBody()
		require.NoError(t, err)
	}, headlessSelect(2))

	assert.Nil(t, body)
}

func TestCollectBodyHeadlessFileSource(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "body.json")
	require.NoError(t, os.WriteFile(filePath, []byte(`{"fileKey": "fileValue"}`), 0o600))

	var body config.HttpBody

	withHeadlessTerminal(t, func() {
		var err error
		body, err = CollectBody()
		require.NoError(t, err)
	}, headlessSelect(0), filePath+headlessEnter)

	assert.Equal(t, config.HttpBody{"fileKey": "fileValue"}, body)
}

func TestCollectBodyStringHeadlessSkipSource(t *testing.T) {
	var bodyString string

	withHeadlessTerminal(t, func() {
		var err error
		bodyString, err = CollectBodyString()
		require.NoError(t, err)
	}, headlessSelect(2))

	assert.Empty(t, bodyString)
}

func TestCollectBodyStringHeadlessFileSource(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "body.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("file contents"), 0o600))

	var bodyString string

	withHeadlessTerminal(t, func() {
		var err error
		bodyString, err = CollectBodyString()
		require.NoError(t, err)
	}, headlessSelect(0), filePath+headlessEnter)

	assert.Equal(t, "file contents", bodyString)
}
