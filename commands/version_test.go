package commands

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsNewer(t *testing.T) {
	assert.True(t, isNewer("v0.5.0", "v0.4.4"))
	assert.True(t, isNewer("0.4.5", "0.4.4"))
	assert.True(t, isNewer("v1.0.0", "v0.9.9"))
	assert.False(t, isNewer("v0.4.4", "v0.4.4"))
	assert.False(t, isNewer("v0.4.3", "v0.4.4"))
	// Tolerates git-describe / pre-release suffixes on the current version.
	assert.False(t, isNewer("v0.4.4", "v0.4.4-2-gabc123"))
	assert.True(t, isNewer("v0.4.5", "v0.4.4-2-gabc123"))
	// Unparseable versions never claim an update.
	assert.False(t, isNewer("nightly", "v0.4.4"))
	assert.False(t, isNewer("v0.5.0", "dev"))
}

func TestParseSemver(t *testing.T) {
	assert.Equal(t, []int{0, 4, 4}, parseSemver("v0.4.4"))
	assert.Equal(t, []int{1, 2, 3}, parseSemver("1.2.3-rc.1"))
	assert.Nil(t, parseSemver("dev"))
	assert.Nil(t, parseSemver("1.2"))
	assert.Nil(t, parseSemver("v1.x.0"))
}

func TestFetchLatestRelease(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v9.9.9"}`))
	}))
	defer srv.Close()
	tag, err := fetchLatestRelease(context.Background(), srv.Client(), srv.URL)
	require.NoError(t, err)
	assert.Equal(t, "v9.9.9", tag)
}

func TestFetchLatestRelease_Errors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	_, err := fetchLatestRelease(context.Background(), srv.Client(), srv.URL)
	assert.Error(t, err)

	// Empty tag is an error.
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv2.Close()
	_, err = fetchLatestRelease(context.Background(), srv2.Client(), srv2.URL)
	assert.Error(t, err)
}
