package auth

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// isolateConfig points the config dir (and therefore the encrypted fallback file
// resolved by NewKeyringStore) at a temp dir, so tests exercising the public API
// never read or write the real ~/.alegra-cli directory.
func isolateConfig(t *testing.T) {
	t.Helper()
	t.Setenv(config.EnvConfig, filepath.Join(t.TempDir(), "config.yaml"))
}

func TestKeyringStore_RoundTrip(t *testing.T) {
	isolateConfig(t)
	keyring.MockInit()
	s := NewKeyringStore()

	require.NoError(t, s.Set("prof", "tok123"))

	got, err := s.Get("prof")
	require.NoError(t, err)
	assert.Equal(t, "tok123", got)
	assert.Equal(t, "tok123", Lookup("prof"))

	require.NoError(t, s.Delete("prof"))
	_, err = s.Get("prof")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Equal(t, "", Lookup("prof"), "Lookup swallows the not-found error")
}

func TestKeyringStore_GetMissing(t *testing.T) {
	isolateConfig(t)
	keyring.MockInit()
	_, err := NewKeyringStore().Get("nope")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestKeyringStore_DeleteMissing(t *testing.T) {
	isolateConfig(t)
	keyring.MockInit()
	err := NewKeyringStore().Delete("nope")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestLookup_KeyringUnavailableFallsBackToFile(t *testing.T) {
	// A locked or unavailable keyring (not just "no entry") must not break auth
	// resolution: with nothing in the fallback file, Lookup returns "" so callers
	// fall back to env/config, and a token written while the keyring is down lands
	// in the encrypted file and reads back transparently.
	isolateConfig(t)
	keyring.MockInitWithError(errors.New("keyring backend locked"))
	t.Cleanup(keyring.MockInit)

	assert.Equal(t, "", Lookup("prof"))
	_, err := NewKeyringStore().Get("prof")
	assert.ErrorIs(t, err, ErrNotFound)

	s := NewKeyringStore()
	require.NoError(t, s.Set("prof", "tok-from-file"))
	assert.Equal(t, "file", s.Backend())
	assert.Equal(t, "tok-from-file", Lookup("prof"))
}
