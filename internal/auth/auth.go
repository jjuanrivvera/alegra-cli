// Package auth stores and retrieves Alegra API tokens. The OS keyring (macOS
// Keychain, Linux Secret Service, Windows Credential Manager) is primary; an
// AES-256-GCM encrypted file is the fallback for headless hosts where no keyring
// backend is available, so tokens are never written to disk in plaintext.
package auth

import "errors"

// keyringService is the service name used in the OS keyring.
const keyringService = "alegra-cli"

// ErrNotFound indicates no token is stored for the profile.
var ErrNotFound = errors.New("no stored token for profile")

// Store persists per-profile API tokens.
type Store interface {
	Set(profile, token string) error
	Get(profile string) (string, error)
	Delete(profile string) error
}

// Lookup tries the keyring (and the encrypted fallback) and returns "" (not an
// error) when nothing is stored or the keyring is unavailable, so callers can fall
// back to config/env.
func Lookup(profile string) string {
	v, err := NewKeyringStore().Get(profile)
	if err != nil {
		return ""
	}
	return v
}
