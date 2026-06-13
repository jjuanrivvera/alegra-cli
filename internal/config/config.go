// Package config loads and persists alegra-cli configuration: named profiles
// (base URL + credentials), global settings, and environment overrides.
//
// Config file: $ALEGRA_CONFIG, else ~/.alegra-cli/config.yaml.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// EnvPrefix-style variable names honored as overrides.
const (
	EnvConfig   = "ALEGRA_CONFIG"
	EnvProfile  = "ALEGRA_PROFILE"
	EnvEmail    = "ALEGRA_EMAIL"
	EnvToken    = "ALEGRA_TOKEN"
	EnvBearer   = "ALEGRA_BEARER_TOKEN"
	EnvBaseURL  = "ALEGRA_BASE_URL"
	EnvOutput   = "ALEGRA_OUTPUT"
	EnvLogLevel = "ALEGRA_LOG_LEVEL"
)

// Profile holds the connection details for one Alegra account.
type Profile struct {
	Name        string `yaml:"-"`
	BaseURL     string `yaml:"baseUrl,omitempty"`
	Email       string `yaml:"email,omitempty"`
	Token       string `yaml:"token,omitempty"`       // optional; prefer keyring
	BearerToken string `yaml:"bearerToken,omitempty"` // optional OAuth
	Description string `yaml:"description,omitempty"`
	// Country caches the account's detected localization (the company's
	// applicationVersion, lowercased — e.g. "colombia", "mexico", "costarica").
	// Auto-populated on `auth login` / refreshed by `doctor`; the platform is the
	// source of truth, so this is a cache, not a user setting.
	Country string `yaml:"country,omitempty"`
}

// Settings holds global, profile-independent options.
type Settings struct {
	DefaultOutputFormat string  `yaml:"defaultOutputFormat,omitempty"`
	RequestsPerSecond   float64 `yaml:"requestsPerSecond,omitempty"`
	LogLevel            string  `yaml:"logLevel,omitempty"`
	// Country is a global, offline fallback hint for pre-flight validation
	// (e.g. "colombia"). The account's real country is auto-detected per profile
	// (Profile.Country); this only applies when nothing has been detected yet.
	Country string `yaml:"country,omitempty"`
}

// Config is the on-disk configuration.
type Config struct {
	DefaultProfile string              `yaml:"defaultProfile,omitempty"`
	Profiles       map[string]*Profile `yaml:"profiles,omitempty"`
	Settings       *Settings           `yaml:"settings,omitempty"`
	// Aliases map a short name to a full command expansion, e.g.
	// "unpaid" -> "invoices list --status open --all".
	Aliases map[string]string `yaml:"aliases,omitempty"`

	path string `yaml:"-"`
}

// DefaultPath returns the resolved config file path.
func DefaultPath() string {
	if p := os.Getenv(EnvConfig); p != "" {
		return p
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".alegra-cli", "config.yaml")
	}
	// HOME unset (stripped CI/cron/container env): fall back to the OS config
	// dir (honors XDG_CONFIG_HOME / %AppData%) so we still resolve an absolute
	// path instead of silently writing config — possibly a token — into the CWD
	// (L10).
	if cfgDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(cfgDir, "alegra-cli", "config.yaml")
	}
	return filepath.Join(".alegra-cli", "config.yaml")
}

// New returns an empty config with defaults populated.
func New() *Config {
	return &Config{
		Profiles: map[string]*Profile{},
		Settings: &Settings{
			DefaultOutputFormat: "table",
			RequestsPerSecond:   5.0,
			LogLevel:            "info",
		},
		path: DefaultPath(),
	}
}

// Load reads the config file, returning a default config if none exists.
func Load() (*Config, error) {
	path := DefaultPath()
	data, err := os.ReadFile(path) //nolint:gosec // path is user-controlled by design
	if err != nil {
		if os.IsNotExist(err) {
			c := New()
			c.path = path
			return c, nil
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	c := New()
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	if c.Profiles == nil {
		c.Profiles = map[string]*Profile{}
	}
	if c.Settings == nil {
		c.Settings = New().Settings
	}
	c.path = path
	return c, nil
}

// Path returns the file path the config was loaded from / will save to.
func (c *Config) Path() string { return c.path }

// Save writes the config to disk (0600), creating the directory as needed.
// The write is atomic (temp + rename) so a crash mid-write can never leave a
// torn config.yaml behind.
func (c *Config) Save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	tmp, err := os.CreateTemp(dir, "config-*.yaml")
	if err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }() // no-op once renamed
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("writing config: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("writing config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	if err := os.Rename(tmp.Name(), c.path); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// ActiveProfileName resolves which profile to use: explicit override > env >
// configured default > "default".
func (c *Config) ActiveProfileName(override string) string {
	if override != "" {
		return override
	}
	if env := os.Getenv(EnvProfile); env != "" {
		return env
	}
	if c.DefaultProfile != "" {
		return c.DefaultProfile
	}
	return "default"
}

// Profile returns the named profile, creating an empty in-memory one if it does
// not exist in the file (so env-only configurations work without a config file).
func (c *Config) Profile(name string) *Profile {
	if p, ok := c.Profiles[name]; ok && p != nil {
		p.Name = name
		return p
	}
	return &Profile{Name: name}
}

// SetProfile inserts or replaces a profile.
func (c *Config) SetProfile(p *Profile) {
	if c.Profiles == nil {
		c.Profiles = map[string]*Profile{}
	}
	c.Profiles[p.Name] = p
}

// Resolved is a fully-resolved connection config after applying env overrides.
type Resolved struct {
	Profile           string
	BaseURL           string
	Email             string
	Token             string
	BearerToken       string
	OutputFormat      string
	RequestsPerSecond float64
	LogLevel          string
}

// Resolve merges a profile with environment overrides and global settings.
// Credential precedence: env > profile file. (Keyring is layered in by the
// caller, which has access to the auth package.)
func (c *Config) Resolve(profileName string) *Resolved {
	p := c.Profile(profileName)
	r := &Resolved{
		Profile:           profileName,
		BaseURL:           firstNonEmpty(os.Getenv(EnvBaseURL), p.BaseURL, "https://api.alegra.com/api/v1"),
		Email:             firstNonEmpty(os.Getenv(EnvEmail), p.Email),
		Token:             firstNonEmpty(os.Getenv(EnvToken), p.Token),
		BearerToken:       firstNonEmpty(os.Getenv(EnvBearer), p.BearerToken),
		OutputFormat:      firstNonEmpty(os.Getenv(EnvOutput), c.Settings.DefaultOutputFormat, "table"),
		RequestsPerSecond: c.Settings.RequestsPerSecond,
		LogLevel:          firstNonEmpty(os.Getenv(EnvLogLevel), c.Settings.LogLevel, "info"),
	}
	if r.RequestsPerSecond <= 0 {
		r.RequestsPerSecond = 5.0
	}
	return r
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// ParseRPS parses a requests-per-second string, returning fallback on error.
func ParseRPS(s string, fallback float64) float64 {
	if s == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f <= 0 {
		return fallback
	}
	return f
}
