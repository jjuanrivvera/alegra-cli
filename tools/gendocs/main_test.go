package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRun generates the command reference into a temp working directory and
// checks the expected pages are written (exercises run()).
func TestRun(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := run(); err != nil {
		t.Fatalf("run() error: %v", err)
	}

	for _, f := range []string{
		filepath.Join(outDir, "index.md"),
		filepath.Join(outDir, "alegra.md"),
		filepath.Join(outDir, "alegra_contacts.md"),
	} {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected generated %s: %v", f, err)
		}
	}
}
