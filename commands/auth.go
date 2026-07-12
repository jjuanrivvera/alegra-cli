package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/auth"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
	"github.com/jjuanrivvera/alegra-cli/internal/version"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Alegra API authentication",
	}
	authCmd.AddCommand(newAuthLoginCmd(), newAuthLogoutCmd(), newAuthStatusCmd())
	rootCmd.AddCommand(authCmd)
}

func newAuthLoginCmd() *cobra.Command {
	var email, token, baseURL string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Store Alegra credentials for a profile (token kept in the OS keyring)",
		Long: `Authenticate against the Alegra API.

Your email and the chosen base URL are saved to the config file; the API token
is stored securely in the OS keyring (not written to disk in plaintext).

Find your API token in Alegra: Configuración → Integraciones → API.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			profile := cfg.ActiveProfileName(flagProfile)
			reader := bufio.NewReader(cmd.InOrStdin())

			if email == "" {
				fmt.Fprint(cmd.OutOrStdout(), "Alegra email: ")
				line, _ := reader.ReadString('\n')
				email = strings.TrimSpace(line)
			}
			if token == "" {
				token, err = promptSecret(cmd, "Alegra API token: ")
				if err != nil {
					return err
				}
			}
			if email == "" || token == "" {
				return fmt.Errorf("email and token are required")
			}

			country, err := loginAndSave(cmd.Context(), cfg, profile, email, token, baseURL)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s (profile %q)\n", email, profile)
			if country != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Detected country (version): %s\n", country)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "Alegra account email")
	cmd.Flags().StringVar(&token, "token", "", "Alegra API token (omit to be prompted securely)")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "Override the API base URL for this profile")
	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove the stored token for the active profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			profile := cfg.ActiveProfileName(flagProfile)
			if err := auth.NewKeyringStore().Delete(profile); err != nil && err != auth.ErrNotFound {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Logged out of profile %q\n", profile)
			return nil
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Aliases: []string{"whoami"},
		Short:   "Show the active profile and verify the stored credentials",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			profile := cfg.ActiveProfileName(flagProfile)
			r := cfg.Resolve(profile)
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Profile:  %s\n", profile)
			fmt.Fprintf(out, "Base URL: %s\n", r.BaseURL)
			fmt.Fprintf(out, "Email:    %s\n", orNone(r.Email))

			var self map[string]any
			if err := client.GetInto(cmd.Context(), "users/self", nil, &self); err != nil {
				fmt.Fprintf(out, "Status:   NOT authenticated (%v)\n", err)
				return err
			}
			fmt.Fprintf(out, "Status:   authenticated ✓\n")
			if name, ok := self["name"]; ok {
				fmt.Fprintf(out, "User:     %v\n", name)
			}
			return nil
		},
	}
}

// loginAndSave verifies credentials against /users/self, detects and caches the
// account's country (applicationVersion), and persists the profile (email +
// base URL to config, token to the keyring). It returns the detected country.
// Shared by `auth login` and `init`.
func loginAndSave(ctx context.Context, cfg *config.Config, profile, email, token, baseURL string) (string, error) {
	base := baseURL
	if base == "" {
		base = cfg.Resolve(profile).BaseURL
	}
	// This is the first time the email+token leave the machine; warn before
	// sending them to a non-HTTPS base URL (M6).
	warnInsecureBaseURL(os.Stderr, base)
	client := api.New(
		api.WithBaseURL(base),
		api.WithBasicAuth(email, token),
		api.WithUserAgent(version.UserAgent()),
	)
	if err := verifyCredentials(ctx, client); err != nil {
		return "", fmt.Errorf("credential check failed: %w", err)
	}

	p := cfg.Profile(profile)
	p.Email = email
	if baseURL != "" {
		p.BaseURL = baseURL
	}
	p.Token = "" // ensure no plaintext token lingers in config
	p.Country = detectCountry(ctx, client)
	cfg.SetProfile(p)
	if cfg.DefaultProfile == "" {
		cfg.DefaultProfile = profile
	}
	if err := cfg.Save(); err != nil {
		return "", err
	}
	if err := auth.NewKeyringStore().Set(profile, token); err != nil {
		return "", fmt.Errorf("storing token in keyring: %w", err)
	}
	return p.Country, nil
}

func verifyCredentials(ctx context.Context, client *api.Client) error {
	var self map[string]any
	return client.GetInto(ctx, "users/self", nil, &self)
}

func promptSecret(cmd *cobra.Command, prompt string) (string, error) {
	fmt.Fprint(cmd.OutOrStdout(), prompt)
	if f, ok := cmd.InOrStdin().(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		s, err := readSecretRaw(f)
		fmt.Fprintln(cmd.OutOrStdout())
		if err != nil {
			return "", err
		}
		return sanitizeSecret(s), nil
	}
	// Non-interactive: read a line.
	line, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// sanitizeSecret strips terminal bracketed-paste markers (ESC[200~ … ESC[201~) and trims
// surrounding whitespace. With bracketed paste enabled, a raw read (unlike the shell's line
// editor) receives those wrappers around pasted text; left in they corrupt a pasted key so it
// fails auth. Stripping them fixes the common "typing works, pasting fails".
func sanitizeSecret(s string) string {
	s = strings.ReplaceAll(s, "\x1b[200~", "")
	s = strings.ReplaceAll(s, "\x1b[201~", "")
	return strings.TrimSpace(s)
}

// readSecretRaw puts the terminal in raw mode (no echo, no line-length limit) and reads one line.
// term.ReadPassword instead reads in CANONICAL mode, capped at MAX_CANON (1024 bytes on macOS):
// pasting a longer secret fills the line buffer and the terminal BLOCKS — the "prompt hangs until
// Ctrl-C" bug. Raw mode has no such limit.
func readSecretRaw(f *os.File) (string, error) {
	fd := int(f.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(fd, oldState) }()
	return scanSecretLine(f)
}

// scanSecretLine reads bytes until CR/LF with no line-length limit. Ctrl-C cancels; Backspace/DEL
// edits. Split out so the byte handling is testable without a real terminal.
func scanSecretLine(r io.Reader) (string, error) {
	var buf []byte
	chunk := make([]byte, 256)
	for {
		n, readErr := r.Read(chunk)
		for i := 0; i < n; i++ {
			switch c := chunk[i]; c {
			case '\r', '\n':
				return string(buf), nil
			case 3: // Ctrl-C
				return "", fmt.Errorf("cancelled")
			case 127, 8: // DEL / Backspace
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
				}
			default:
				buf = append(buf, c)
			}
		}
		if readErr != nil {
			if len(buf) == 0 {
				return "", readErr
			}
			return string(buf), nil
		}
	}
}

func orNone(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}
