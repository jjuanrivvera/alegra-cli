package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/auth"
	"github.com/jjuanrivvera/alegra-cli/internal/catalog"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
	"github.com/jjuanrivvera/alegra-cli/internal/ui"
)

func init() {
	rootCmd.AddCommand(readOnlyHints(&cobra.Command{
		Use:   "doctor",
		Short: "Diagnose configuration, credentials, and account health",
		Long: `doctor runs a series of read-only checks: configuration, credentials,
authentication, company/country, rate-limit budget, numbering resolutions, and a
plan probe. Use it first whenever something isn't working.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDoctor(cmd)
		},
	}))
}

func runDoctor(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()
	c := ui.For(out)
	ok := func(label, detail string) { fmt.Fprintf(out, "%s %-13s %s\n", c.Green("✔"), label, detail) }
	warn := func(label, detail string) { fmt.Fprintf(out, "%s %-13s %s\n", c.Yellow("⚠"), label, detail) }
	bad := func(label, detail string) { fmt.Fprintf(out, "%s %-13s %s\n", c.Red("✘"), label, detail) }

	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		bad("config", err.Error())
		return err
	}
	profile := cfg.ActiveProfileName(flagProfile)
	ok("config", fmt.Sprintf("%s (profile: %s)", cfg.Path(), profile))

	// 2. Credentials
	r := cfg.Resolve(profile)
	tokenSrc := ""
	switch {
	case r.BearerToken != "":
		tokenSrc = "bearer token"
	case r.Token != "":
		tokenSrc = "config/env token"
	case auth.Lookup(profile) != "":
		tokenSrc = "keyring token"
	}
	if tokenSrc == "" {
		bad("credentials", "none found — run `alegra auth login` or set ALEGRA_EMAIL/ALEGRA_TOKEN")
		return fmt.Errorf("no credentials for profile %q", profile)
	}
	ok("credentials", tokenSrc)

	client, err := getAPIClient(cmd)
	if err != nil {
		bad("client", err.Error())
		return err
	}

	// 3. Auth
	var self map[string]any
	if err := client.GetInto(cmd.Context(), "users/self", nil, &self); err != nil {
		bad("auth", err.Error())
		return err
	}
	ok("auth", fmt.Sprintf("%v (%v)", strOr(self, "name", "?"), strOr(self, "role", "user")))

	// 4. Company / country
	var company map[string]any
	if err := client.GetInto(cmd.Context(), "company", nil, &company); err != nil {
		warn("company", err.Error())
	} else {
		ok("company", fmt.Sprintf("%v · country: %v · regime: %v",
			strOr(company, "name", "?"), strOr(company, "applicationVersion", "?"), strOr(company, "regime", "—")))
		// Refresh the cached per-profile country used for pre-flight validation.
		if v, isStr := company["applicationVersion"].(string); isStr {
			cacheCountry(cfg, profile, strings.ToLower(strings.TrimSpace(v)))
		}
		// Mexican accounts rely on the locally-synced SAT product-keys catalog
		// (no Alegra endpoint exists); surface its state where users debug.
		if strings.EqualFold(strOr(company, "applicationVersion", ""), "mexico") {
			if dir, derr := catalogsDir(); derr == nil {
				if sat, lerr := catalog.LoadSAT(dir); lerr == nil {
					ok("sat catalog", fmt.Sprintf("%d product keys (version %s, synced %s)",
						len(sat.Entries), sat.Version, sat.FetchedAt.Format("2006-01-02")))
				} else {
					warn("sat catalog", "not synced — run `alegra catalog sync-sat` for offline product-key search")
				}
			}
		}
	}

	// 5. Rate-limit budget (populated by the calls above)
	if limit, remaining, reset := client.RateLimit(); limit > 0 {
		ok("rate limit", fmt.Sprintf("%d/%d remaining (resets in %ds)", remaining, limit, reset))
	} else {
		warn("rate limit", "no X-Rate-Limit headers seen")
	}

	// 6. Numbering resolutions
	templates, err := client.NumberTemplates().List(cmd.Context(), api.ListParams{})
	if err != nil {
		warn("numbering", err.Error())
	} else {
		ok("numbering", fmt.Sprintf("%d resolution(s) configured", len(templates)))
	}

	// 7. Plan probe (a commonly plan-gated endpoint). The documented path is
	// /conciliations (no "re-"); probing the wrong path returns a misleading 403.
	if err := client.GetInto(cmd.Context(), "conciliations", nil, &[]any{}); err != nil {
		if apiErr, isAPI := api.AsAPIError(err); isAPI && (apiErr.StatusCode == 402 || apiErr.StatusCode == 403) {
			warn("plan", fmt.Sprintf("'conciliations' returned %d — not in your plan", apiErr.StatusCode))
		}
		// other errors here are not meaningful for a health check
	} else {
		ok("plan", "conciliations accessible")
	}

	fmt.Fprintln(out, "\n"+c.Green("All critical checks passed."))
	return nil
}

func strOr(m map[string]any, key, def string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return def
}
