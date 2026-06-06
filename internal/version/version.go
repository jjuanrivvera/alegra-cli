// Package version exposes build metadata injected at link time via -ldflags.
package version

import (
	"fmt"
	"runtime"
)

// These are set via -ldflags at build time (see Makefile / .goreleaser.yaml).
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info returns a one-line human-readable version string.
func Info() string {
	return fmt.Sprintf("alegra %s (commit %s, built %s, %s/%s, %s)",
		Version, Commit, BuildDate, runtime.GOOS, runtime.GOARCH, runtime.Version())
}

// Short returns just the semantic version.
func Short() string { return Version }

// UserAgent returns the HTTP User-Agent string for API requests.
func UserAgent() string {
	return fmt.Sprintf("alegra-cli/%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
}
