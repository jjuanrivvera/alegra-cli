package commands

import (
	"bufio"
	"context"
	"fmt"
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

			// Verify the credentials against /users/self before saving.
			base := baseURL
			if base == "" {
				base = cfg.Resolve(profile).BaseURL
			}
			client := api.New(
				api.WithBaseURL(base),
				api.WithBasicAuth(email, token),
				api.WithUserAgent(version.UserAgent()),
			)
			if err := verifyCredentials(cmd.Context(), client); err != nil {
				return fmt.Errorf("credential check failed: %w", err)
			}

			// Persist: email + base URL in config, token in keyring.
			p := cfg.Profile(profile)
			p.Email = email
			if baseURL != "" {
				p.BaseURL = baseURL
			}
			p.Token = "" // ensure no plaintext token lingers in config
			cfg.SetProfile(p)
			if cfg.DefaultProfile == "" {
				cfg.DefaultProfile = profile
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			if err := auth.NewKeyringStore().Set(profile, token); err != nil {
				return fmt.Errorf("storing token in keyring: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s (profile %q)\n", email, profile)
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
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			cfg, _ := config.Load()
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

func verifyCredentials(ctx context.Context, client *api.Client) error {
	var self map[string]any
	return client.GetInto(ctx, "users/self", nil, &self)
}

func promptSecret(cmd *cobra.Command, prompt string) (string, error) {
	fmt.Fprint(cmd.OutOrStdout(), prompt)
	if f, ok := cmd.InOrStdin().(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		b, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(cmd.OutOrStdout())
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	// Non-interactive: read a line.
	line, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func orNone(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}
