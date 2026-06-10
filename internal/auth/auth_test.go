package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestKeyringStore_RoundTrip(t *testing.T) {
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
	keyring.MockInit()
	_, err := NewKeyringStore().Get("nope")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestKeyringStore_DeleteMissing(t *testing.T) {
	keyring.MockInit()
	err := NewKeyringStore().Delete("nope")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestLookup_KeyringBackendFailureDegradesGracefully(t *testing.T) {
	// A locked or unavailable keyring (not just "no entry") must not break
	// auth resolution: Lookup returns "" so callers fall back to env/config.
	keyring.MockInitWithError(errors.New("keyring backend locked"))
	t.Cleanup(keyring.MockInit)

	assert.Equal(t, "", Lookup("prof"))

	// The raw store surfaces the backend error untouched for callers that
	// need to distinguish it from not-found.
	_, err := NewKeyringStore().Get("prof")
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNotFound)
}
