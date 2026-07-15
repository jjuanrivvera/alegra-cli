package auth

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/zalando/go-keyring"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// KeyringStore persists per-profile API tokens in the OS keyring (macOS Keychain,
// Linux Secret Service, Windows Credential Manager), transparently falling back to
// an encrypted file when the keyring is unavailable — a headless Linux box with no
// Secret Service, for example. Callers see one Store; the backend switch is invisible.
type KeyringStore struct {
	service string
	fb      *fileStore
	backend string
}

// Compile-time check that KeyringStore satisfies the public Store contract.
var _ Store = (*KeyringStore)(nil)

// NewKeyringStore returns the default token store. The encrypted fallback file lives
// next to config.yaml (the config dir) and is only touched when the keyring is
// unreachable.
func NewKeyringStore() *KeyringStore { return newStore(defaultFallbackDir()) }

// newStore builds a store whose fallback file lives under dir. Tests use it with a
// temp dir so they never touch the real config directory or a live OS keyring.
func newStore(dir string) *KeyringStore {
	return &KeyringStore{service: keyringService, fb: newFileStore(dir), backend: "keyring"}
}

// defaultFallbackDir is the directory holding config.yaml, so credentials.enc sits
// beside it and honors $ALEGRA_CONFIG / XDG overrides exactly as the config does.
func defaultFallbackDir() string { return filepath.Dir(config.DefaultPath()) }

// Backend reports where the token currently lives ("keyring" or "file"), for doctor output.
func (s *KeyringStore) Backend() string { return s.backend }

// Set stores the token for a profile, falling back to the encrypted file on any
// keyring write error.
func (s *KeyringStore) Set(profile, token string) error {
	if err := keyring.Set(s.service, profile, token); err != nil {
		s.backend = "file"
		return s.fb.Set(profile, token)
	}
	return nil
}

// Get retrieves the token for a profile. It falls back to the encrypted file both
// when the keyring is reachable but empty (the token may have been written on a
// headless host) and when the keyring is entirely unavailable.
func (s *KeyringStore) Get(profile string) (string, error) {
	tok, err := keyring.Get(s.service, profile)
	if err == nil {
		return tok, nil
	}
	if errors.Is(err, keyring.ErrNotFound) {
		// Keyring works but has nothing — check the fallback file before giving up.
		if tok, ferr := s.fb.Get(profile); ferr == nil {
			s.backend = "file"
			return tok, nil
		}
		return "", ErrNotFound
	}
	// Keyring is unavailable entirely → use the fallback file.
	s.backend = "file"
	tok, ferr := s.fb.Get(profile)
	if ferr != nil {
		return "", ErrNotFound
	}
	return tok, nil
}

// Delete removes the stored token from both backends. The keyring delete is
// best-effort (it may be entirely unavailable, in which case the encrypted file is
// the real store). Deleting a profile absent from both backends returns ErrNotFound,
// preserving the store's original contract.
func (s *KeyringStore) Delete(profile string) error {
	if kerr := keyring.Delete(s.service, profile); kerr == nil {
		// Also drop any file copy so the two backends can't diverge; best effort.
		_ = s.fb.Delete(profile)
		return nil
	}
	// Keyring had nothing (or is unavailable) — the file is the source of truth.
	ferr := s.fb.Delete(profile)
	if ferr == nil || errors.Is(ferr, ErrNotFound) {
		return ferr
	}
	return fmt.Errorf("delete token: %w", ferr)
}
