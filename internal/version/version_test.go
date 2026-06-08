package version

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShort(t *testing.T) {
	old := Version
	defer func() { Version = old }()
	Version = "v9.9.9"
	assert.Equal(t, "v9.9.9", Short())
}

func TestInfoContainsBuildMetadata(t *testing.T) {
	oV, oC, oB := Version, Commit, BuildDate
	defer func() { Version, Commit, BuildDate = oV, oC, oB }()
	Version, Commit, BuildDate = "v1.2.3", "abc123", "2026-06-08"

	got := Info()
	for _, want := range []string{
		"alegra", "v1.2.3", "abc123", "2026-06-08",
		runtime.GOOS, runtime.GOARCH, runtime.Version(),
	} {
		assert.Contains(t, got, want)
	}
}

func TestUserAgent(t *testing.T) {
	old := Version
	defer func() { Version = old }()
	Version = "v1.2.3"

	got := UserAgent()
	assert.True(t, strings.HasPrefix(got, "alegra-cli/v1.2.3"), "got %q", got)
	assert.Contains(t, got, runtime.GOOS)
	assert.Contains(t, got, runtime.GOARCH)
}
