package auth

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	fs := newFileStore(dir)

	require.NoError(t, fs.Set("a", "token-a"))
	require.NoError(t, fs.Set("b", "token-b"))

	got, err := fs.Get("a")
	require.NoError(t, err)
	assert.Equal(t, "token-a", got)

	// The on-disk file must be 0600 (Unix only — Windows has no POSIX mode bits) and must
	// NOT contain the plaintext token.
	info, err := os.Stat(filepath.Join(dir, "credentials.enc"))
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "credentials.enc"))
	assert.NotContains(t, string(raw), "token-a", "token must be encrypted at rest")

	require.NoError(t, fs.Delete("a"))
	_, err = fs.Get("a")
	assert.ErrorIs(t, err, ErrNotFound)

	// b survives a's deletion.
	got, err = fs.Get("b")
	require.NoError(t, err)
	assert.Equal(t, "token-b", got)
}

func TestFileStore_PasswordKey(t *testing.T) {
	t.Setenv("ALEGRA_KEYRING_PASSWORD", "correct horse battery staple")
	dir := t.TempDir()
	fs := newFileStore(dir)
	require.NoError(t, fs.Set("x", "secret"))

	got, err := fs.Get("x")
	require.NoError(t, err)
	assert.Equal(t, "secret", got)

	// A different password must fail to decrypt.
	t.Setenv("ALEGRA_KEYRING_PASSWORD", "wrong password entirely")
	_, err = fs.Get("x")
	assert.Error(t, err)
}

func TestEncryptDecrypt(t *testing.T) {
	ct, err := encrypt("hello")
	require.NoError(t, err)
	pt, err := decrypt(ct)
	require.NoError(t, err)
	assert.Equal(t, "hello", pt)

	_, err = decrypt("not-base64!!!")
	assert.Error(t, err)
	_, err = decrypt("dG9vc2hvcnQ=") // valid base64 but too short for a nonce
	assert.Error(t, err)
}
