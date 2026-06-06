// Package auth stores and retrieves Alegra API tokens in the OS keyring, with
// graceful fallback when no keyring backend is available.
package auth

import (
	"errors"

	"github.com/zalando/go-keyring"
)

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

// KeyringStore uses the OS keyring (macOS Keychain, Linux Secret Service,
// Windows Credential Manager).
type KeyringStore struct{}

// NewKeyringStore returns a keyring-backed store.
func NewKeyringStore() *KeyringStore { return &KeyringStore{} }

// Set stores the token for a profile.
func (s *KeyringStore) Set(profile, token string) error {
	return keyring.Set(keyringService, profile, token)
}

// Get retrieves the token for a profile.
func (s *KeyringStore) Get(profile string) (string, error) {
	v, err := keyring.Get(keyringService, profile)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return v, nil
}

// Delete removes the stored token for a profile.
func (s *KeyringStore) Delete(profile string) error {
	err := keyring.Delete(keyringService, profile)
	if err != nil && errors.Is(err, keyring.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// Lookup tries the keyring and returns "" (not an error) when nothing is stored
// or the keyring is unavailable, so callers can fall back to config/env.
func Lookup(profile string) string {
	v, err := NewKeyringStore().Get(profile)
	if err != nil {
		return ""
	}
	return v
}
