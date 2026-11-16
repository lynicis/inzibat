package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Key sequences for driving huh forms through a swapped stdin; selects only
// accept arrow keys, plain text would land in the select's filter and panic
// on a non-TTY output. Confirm fields advance on the bare y/n key.
const (
	formbuilderHeadlessEnter = "\r"
	formbuilderHeadlessYes   = "y"
	formbuilderHeadlessNo    = "n"
)

// formbuilderHeadlessSelect picks the zero-based index-th option.
func formbuilderHeadlessSelect(index int) string {
	return strings.Repeat("\x1b[B", index) + formbuilderHeadlessEnter
}

// withHeadlessTerminal runs fn with os.Stdin replaced by a pipe fed the key
// segments and os.Stdout sent to /dev/null, so real huh forms run headless.
// TERM is forced to a real value: huh silently switches to accessible mode
// under TERM=dumb, which never populates form results. Each segment is one
// field's keys, written separately: a form drops the enter between two
// fields when a whole multi-field stream arrives in a single read, and a
// previous form's reader slurps segments written before it exits.
// Tests using it must not run in parallel.
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

func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}
