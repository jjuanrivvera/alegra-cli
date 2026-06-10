package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/catalog"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
	"github.com/jjuanrivvera/alegra-cli/internal/ui"
)

func init() {
	var email, token, baseURL string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Guided setup: authenticate, detect your country, and save a profile",
		Long: `init walks you through first-time setup: it takes your Alegra email and
API token, verifies them, auto-detects your account's country (the API version),
and saves a profile (token in the OS keyring).

Find your API token in Alegra: Configuración → Integraciones → API.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runInit(cmd, email, token, baseURL)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "Alegra account email (skips the prompt)")
	cmd.Flags().StringVar(&token, "token", "", "Alegra API token (skips the prompt)")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "Override the API base URL")
	rootCmd.AddCommand(cmd)
}

func runInit(cmd *cobra.Command, email, token, baseURL string) error {
	out := cmd.OutOrStdout()
	c := ui.For(out)

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	profile := cfg.ActiveProfileName(flagProfile)

	fmt.Fprintln(out, c.Bold("Welcome to alegra-cli")+" — let's get you set up.")
	fmt.Fprintf(out, "%s\n\n", c.Dim(fmt.Sprintf("Configuring profile %q (config: %s)", profile, cfg.Path())))

	reader := bufio.NewReader(cmd.InOrStdin())
	if email == "" {
		email = strings.TrimSpace(os.Getenv(config.EnvEmail))
	}
	if email == "" {
		fmt.Fprint(out, "Alegra email: ")
		line, _ := reader.ReadString('\n')
		email = strings.TrimSpace(line)
	}
	if token == "" {
		token = strings.TrimSpace(os.Getenv(config.EnvToken))
	}
	if token == "" {
		token, err = promptSecret(cmd, "Alegra API token (input hidden): ")
		if err != nil {
			return err
		}
	}
	if email == "" || token == "" {
		return fmt.Errorf("email and token are required")
	}

	fmt.Fprint(out, "\nVerifying credentials… ")
	country, err := loginAndSave(cmd.Context(), cfg, profile, email, token, baseURL)
	if err != nil {
		fmt.Fprintln(out, c.Red("failed"))
		return err
	}
	fmt.Fprintln(out, c.Green("ok"))

	fmt.Fprintf(out, "\n%s\n", c.Green("✔ Setup complete."))
	fmt.Fprintf(out, "  %s %s\n", c.Dim("Account:"), email)
	if country != "" {
		fmt.Fprintf(out, "  %s %s\n", c.Dim("Country (API version):"), country)
	}

	offerSATCatalog(cmd, reader, country)

	fmt.Fprintf(out, "\n%s\n", c.Bold("Next steps:"))
	fmt.Fprintln(out, "  alegra doctor                 # verify auth, plan, rate limit, numbering")
	fmt.Fprintln(out, "  alegra contacts list          # list your clients/suppliers")
	fmt.Fprintln(out, "  alegra invoices list --count  # how many invoices you have")
	fmt.Fprintln(out, "  alegra <resource> --help      # discover any resource's actions")
	return nil
}

// offerSATCatalog proposes the one-time SAT product-keys download when init
// just detected a Mexican account — the moment the user most needs it and
// least knows the command exists. It must never fail init: a declined prompt,
// EOF (non-interactive stdin), or a download error all degrade to a hint.
func offerSATCatalog(cmd *cobra.Command, reader *bufio.Reader, country string) {
	if !strings.EqualFold(strings.TrimSpace(country), "mexico") {
		return
	}
	out := cmd.OutOrStdout()
	c := ui.For(out)
	dir, err := catalogsDir()
	if err != nil || catalog.SATCached(dir) {
		return
	}

	hint := func() {
		fmt.Fprintf(out, "  %s\n", c.Dim("Skipped — run `alegra catalog sync-sat` anytime."))
	}
	fmt.Fprint(out, "\nMexican account detected. Download the SAT product-keys catalog\n(claveProdServ, ~7MB one time, used for offline search and CFDI validation)? [Y/n]: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(out)
		hint()
		return
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "", "y", "yes", "s", "si", "sí":
	default:
		hint()
		return
	}

	fmt.Fprint(out, "Downloading… ")
	cat, err := catalog.SyncSAT(cmd.Context(), dir)
	if err != nil {
		fmt.Fprintln(out, c.Yellow("failed: "+err.Error()))
		hint()
		return
	}
	fmt.Fprintln(out, c.Green(fmt.Sprintf("ok (%d keys, catalog %s)", len(cat.Entries), cat.Version)))
	fmt.Fprintf(out, "  %s\n", c.Dim("Try: alegra catalog product-keys \"refrigerador\""))
}
