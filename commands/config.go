package commands

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage alegra-cli configuration and profiles",
	}
	configCmd.AddCommand(
		newConfigPathCmd(),
		newConfigViewCmd(),
		newConfigSetProfileCmd(),
		newConfigUseCmd(),
		newConfigListProfilesCmd(),
	)
	rootCmd.AddCommand(configCmd)
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), config.DefaultPath())
		},
	}
}

func newConfigViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Show the current configuration (tokens redacted)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "config file:     %s\n", cfg.Path())
			fmt.Fprintf(out, "default profile: %s\n", orNone(cfg.DefaultProfile))
			fmt.Fprintf(out, "output format:   %s\n", cfg.Settings.DefaultOutputFormat)
			fmt.Fprintf(out, "requests/sec:    %.1f\n", cfg.Settings.RequestsPerSecond)
			fmt.Fprintln(out, "profiles:")
			names := make([]string, 0, len(cfg.Profiles))
			for n := range cfg.Profiles {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				p := cfg.Profiles[n]
				tok := "(keyring)"
				if p.Token != "" {
					tok = "(in config, redacted)"
				}
				fmt.Fprintf(out, "  - %s: email=%s baseUrl=%s token=%s\n",
					n, orNone(p.Email), orNone(p.BaseURL), tok)
			}
			return nil
		},
	}
}

func newConfigSetProfileCmd() *cobra.Command {
	var name, email, baseURL string
	cmd := &cobra.Command{
		Use:   "set-profile",
		Short: "Create or update a profile (use `alegra auth login` to set the token)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if name == "" {
				name = cfg.ActiveProfileName(flagProfile)
			}
			p := cfg.Profile(name)
			if email != "" {
				p.Email = email
			}
			if baseURL != "" {
				p.BaseURL = baseURL
			}
			cfg.SetProfile(p)
			if cfg.DefaultProfile == "" {
				cfg.DefaultProfile = name
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved profile %q\n", name)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Profile name (default: active profile)")
	cmd.Flags().StringVar(&email, "email", "", "Account email")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "API base URL")
	return cmd
}

func newConfigUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile>",
		Short: "Set the default profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cfg.DefaultProfile = args[0]
			if _, ok := cfg.Profiles[args[0]]; !ok {
				cfg.SetProfile(&config.Profile{Name: args[0]})
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Default profile set to %q\n", args[0])
			return nil
		},
	}
}

func newConfigListProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list-profiles",
		Aliases: []string{"profiles"},
		Short:   "List configured profiles",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			names := make([]string, 0, len(cfg.Profiles))
			for n := range cfg.Profiles {
				names = append(names, n)
			}
			sort.Strings(names)
			out := cmd.OutOrStdout()
			for _, n := range names {
				marker := "  "
				if n == cfg.DefaultProfile {
					marker = "* "
				}
				fmt.Fprintf(out, "%s%s\n", marker, n)
			}
			if len(names) == 0 {
				fmt.Fprintln(out, "(no profiles configured)")
			}
			return nil
		},
	}
}
