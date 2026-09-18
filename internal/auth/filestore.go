package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
)

// fileStore is the encrypted-file fallback used when no OS keyring is available. Tokens are
// AES-256-GCM encrypted in a 0600 file. The encryption key comes from $ALEGRA_KEYRING_PASSWORD
// (scrypt-derived) when set; otherwise a host-bound key derived from hostname+uid. The
// host-bound key is obfuscation-grade, NOT a security boundary — the OS keyring is the real
// protection, and this fallback only exists so headless hosts still avoid plaintext tokens.
type fileStore struct {
	path string
}

func newFileStore(dir string) *fileStore {
	return &fileStore{path: filepath.Join(dir, "credentials.enc")}
}

// envelope is the on-disk format: profile -> base64(nonce || ciphertext).
type envelope map[string]string

func (f *fileStore) load() (envelope, error) {
	data, err := os.ReadFile(f.path) //nolint:gosec // G304: fixed path under the config dir
	if os.IsNotExist(err) {
		return envelope{}, nil
	}
	if err != nil {
		return nil, err
	}
	env := envelope{}
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse credentials file: %w", err)
	}
	return env, nil
}

func (f *fileStore) save(env envelope) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(f.path), ".cred-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, f.path)
}

func (f *fileStore) Set(profile, token string) error {
	env, err := f.load()
	if err != nil {
		return err
	}
	ct, err := encrypt(token)
	if err != nil {
		return err
	}
	env[profile] = ct
	return f.save(env)
}

func (f *fileStore) Get(profile string) (string, error) {
	env, err := f.load()
	if err != nil {
		return "", err
	}
	ct, ok := env[profile]
	if !ok {
		return "", ErrNotFound
	}
	return decrypt(ct)
}

func (f *fileStore) Delete(profile string) error {
	env, err := f.load()
	if err != nil {
		return err
	}
	if _, ok := env[profile]; !ok {
		return ErrNotFound
	}
	delete(env, profile)
	return f.save(env)
}

// --- crypto ---

func deriveKey() ([]byte, error) {
	salt := []byte("alegra-credentials-v1")
	if pw := os.Getenv("ALEGRA_KEYRING_PASSWORD"); pw != "" {
		return scrypt.Key([]byte(pw), salt, 1<<15, 8, 1, 32)
	}
	host, _ := os.Hostname()
	seed := fmt.Sprintf("alegra-fallback|%s|%d", host, os.Getuid())
	sum := sha256.Sum256(append(salt, seed...))
	return sum[:], nil
}

func gcm() (cipher.AEAD, error) {
	key, err := deriveKey()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func encrypt(plain string) (string, error) {
	aead, err := gcm()
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := aead.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func decrypt(b64 string) (string, error) {
	aead, err := gcm()
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	ns := aead.NonceSize()
	if len(data) < ns {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ct := data[:ns], data[ns:]
	plain, err := aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed (wrong $ALEGRA_KEYRING_PASSWORD or host changed): %w", err)
	}
	return string(plain), nil
}
