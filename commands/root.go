// Package commands implements the alegra CLI command tree. Each Alegra resource
// registers itself via init() (see registerResource), so adding a resource never
// requires editing this file.
package commands

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/auth"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
	"github.com/jjuanrivvera/alegra-cli/internal/output"
	"github.com/jjuanrivvera/alegra-cli/internal/ui"
	"github.com/jjuanrivvera/alegra-cli/internal/version"
)

// Global persistent flags.
var (
	flagProfile   string
	flagOutput    string
	flagBaseURL   string
	flagRPS       float64
	flagDryRun    bool
	flagShowToken bool
	flagVerbose   bool
	flagColumns   []string
	flagNoColor   bool
)

// rootCmd is the top-level command. It is a package var so resource files can
// attach themselves from init() (package vars initialize before init() runs).
var rootCmd = &cobra.Command{
	Use:   "alegra",
	Short: "Alegra accounting system CLI",
	Long: `alegra is a command-line interface for the Alegra accounting API
(https://developer.alegra.com).

It manages contacts, invoices, items, payments, taxes, and the full Alegra
resource surface, with table/json/yaml/csv output, profiles, and a dry-run mode.

Authenticate with:
  alegra auth login            # stores your API token in the OS keyring
or set ALEGRA_EMAIL and ALEGRA_TOKEN in the environment.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&flagProfile, "profile", "", "Configuration profile to use (env: ALEGRA_PROFILE)")
	pf.StringVarP(&flagOutput, "output", "o", "", "Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)")
	pf.StringVar(&flagBaseURL, "base-url", "", "Override the API base URL (env: ALEGRA_BASE_URL)")
	pf.Float64Var(&flagRPS, "requests-per-second", 0, "Client-side rate limit (default from config)")
	pf.BoolVar(&flagDryRun, "dry-run", false, "Print the equivalent curl request without sending it")
	pf.BoolVar(&flagShowToken, "show-token", false, "In --dry-run, do not redact the Authorization header")
	pf.BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose (debug) logging to stderr")
	pf.StringSliceVar(&flagColumns, "columns", nil, "Comma-separated columns for table/csv output")
	pf.BoolVar(&flagNoColor, "no-color", false, "Disable colored output (also respects the NO_COLOR env var)")

	// Dynamic completion for the global flags. --columns is shared by every
	// command, so a single completer resolves the per-resource field set from the
	// target command's annotation (see withColumns/completeColumns).
	_ = rootCmd.RegisterFlagCompletionFunc("output", fixedCompleter("table", "json", "yaml", "csv"))
	_ = rootCmd.RegisterFlagCompletionFunc("profile", completeProfiles)
	_ = rootCmd.RegisterFlagCompletionFunc("columns", completeColumns)

	// Apply the color preference before any command renders output.
	rootCmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		ui.SetNoColor(flagNoColor)
	}
}

// Execute runs the root command. Resources are registered via their init().
func Execute() error {
	// Expand a leading user alias (if any) before cobra parses the arguments.
	if cfg, err := config.Load(); err == nil && len(cfg.Aliases) > 0 {
		os.Args = expandAliasArgs(os.Args, cfg.Aliases, isBuiltinCommand)
	}
	return rootCmd.Execute()
}

// RootCmd exposes the root command (used by docs generation and the MCP bridge).
func RootCmd() *cobra.Command { return rootCmd }

// --- shared configuration resolution (loaded once) ---

var (
	cfgOnce  sync.Once
	cfg      *config.Config
	resolved *config.Resolved
	cfgErr   error
)

func loadResolved() (*config.Config, *config.Resolved, error) {
	cfgOnce.Do(func() {
		cfg, cfgErr = config.Load()
		if cfgErr != nil {
			return
		}
		name := cfg.ActiveProfileName(flagProfile)
		resolved = cfg.Resolve(name)
		// Flag overrides take precedence over config/env.
		if flagBaseURL != "" {
			resolved.BaseURL = flagBaseURL
		}
		if flagOutput != "" {
			resolved.OutputFormat = flagOutput
		}
		if flagRPS > 0 {
			resolved.RequestsPerSecond = flagRPS
		}
		// Layer in keyring credentials if the profile has no token yet.
		if resolved.Token == "" && resolved.BearerToken == "" {
			if tok := auth.Lookup(name); tok != "" {
				resolved.Token = tok
			}
		}
	})
	return cfg, resolved, cfgErr
}

// getAPIClient builds an API client from the resolved configuration and flags.
// cmd supplies the dry-run output sink (cmd.OutOrStdout()) so the printed curl
// honors any output redirection and is capturable in tests.
func getAPIClient(cmd *cobra.Command) (*api.Client, error) {
	_, r, err := loadResolved()
	if err != nil {
		return nil, err
	}

	if !flagDryRun && r.Token == "" && r.BearerToken == "" {
		return nil, fmt.Errorf("no credentials for profile %q: run `alegra auth login` or set ALEGRA_EMAIL and ALEGRA_TOKEN", r.Profile)
	}

	if r.Token != "" || r.BearerToken != "" {
		warnInsecureBaseURL(os.Stderr, r.BaseURL)
	}

	logger := newLogger(r.LogLevel)
	opts := []api.Option{
		api.WithBaseURL(r.BaseURL),
		api.WithRequestsPerSecond(r.RequestsPerSecond),
		api.WithUserAgent(version.UserAgent()),
		api.WithLogger(logger),
		api.WithDryRun(flagDryRun, flagShowToken, cmd.OutOrStdout()),
	}
	if r.BearerToken != "" {
		opts = append(opts, api.WithBearerToken(r.BearerToken))
	} else {
		opts = append(opts, api.WithBasicAuth(r.Email, r.Token))
	}
	return api.New(opts...), nil
}

// warnInsecureBaseURL prints a warning when API credentials would be sent to a
// non-HTTPS, non-loopback base URL (cleartext credential exposure). Only https,
// and http to loopback (local sandboxes/tests), are silently safe. A schemeless
// URL (e.g. "api.alegra.com/api/v1") is warned about too: it is not https and
// will produce a confusing low-level request error otherwise (L12).
func warnInsecureBaseURL(w io.Writer, baseURL string) {
	u, err := url.Parse(baseURL)
	if err != nil || u.Scheme == "https" {
		return
	}
	host := u.Hostname()
	if host == "localhost" {
		return
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return
	}
	fmt.Fprintf(w, "warning: base URL %q is not HTTPS — API credentials will be sent in cleartext\n", baseURL)
}

func newLogger(level string) *slog.Logger {
	lvl := slog.LevelInfo
	if flagVerbose {
		lvl = slog.LevelDebug
	} else {
		switch strings.ToLower(level) {
		case "debug":
			lvl = slog.LevelDebug
		case "warn", "warning":
			lvl = slog.LevelWarn
		case "error":
			lvl = slog.LevelError
		}
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}

// outputFormat returns the effective output format.
func outputFormat() output.Format {
	_, r, err := loadResolved()
	if err != nil || r == nil {
		if flagOutput != "" {
			return output.Format(strings.ToLower(flagOutput))
		}
		return output.FormatTable
	}
	return output.Format(strings.ToLower(r.OutputFormat))
}

// render writes data using the resolved format and column selection.
func render(cmd *cobra.Command, data any, defaultCols []string) error {
	cols := flagColumns
	if len(cols) == 0 {
		cols = defaultCols
	}
	return output.Render(cmd.OutOrStdout(), data, outputFormat(), cols)
}
