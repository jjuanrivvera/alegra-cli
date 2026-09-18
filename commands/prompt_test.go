package commands

import (
	"strings"
	"testing"
)

func TestScanSecretLine(t *testing.T) {
	// A long paste (well past canonical mode's 1024-char MAX_CANON) reads intact — the whole point.
	long := strings.Repeat("A", 2000)
	if got, err := scanSecretLine(strings.NewReader(long + "\r")); err != nil || got != long {
		t.Errorf("long line: len=%d err=%v", len(got), err)
	}
	// Stops at LF; Backspace (0x7f) edits.
	if got, _ := scanSecretLine(strings.NewReader("ab\x7fc\n")); got != "ac" {
		t.Errorf("backspace: %q", got)
	}
	// Ctrl-C cancels.
	if _, err := scanSecretLine(strings.NewReader("\x03")); err == nil {
		t.Error("Ctrl-C should cancel")
	}
	// EOF with buffered content still returns it.
	if got, _ := scanSecretLine(strings.NewReader("xyz")); got != "xyz" {
		t.Errorf("EOF: %q", got)
	}
}

func TestSanitizeSecret(t *testing.T) {
	key := "eyJhbGci.payload.sig"
	if got := sanitizeSecret("\x1b[200~" + key + "\x1b[201~\n"); got != key {
		t.Errorf("bracketed paste not stripped: %q", got)
	}
	if got := sanitizeSecret("  " + key + "  "); got != key {
		t.Errorf("trim: %q", got)
	}
}
