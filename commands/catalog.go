package commands

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/catalog"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// catalogCategoryRow is the table row for `alegra catalog` (no args).
type catalogCategoryRow struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Entries  int    `json:"entries"`
}

func init() {
	var country string
	cmd := &cobra.Command{
		Use:     "catalog [category]",
		Aliases: []string{"catalogs", "reference"},
		Short:   "Country reference catalogs (units, identification types, taxes, ...)",
		Long: "Show Alegra's per-country reference catalogs — units of measure and reference\n" +
			"enums such as identification types, tax types, payment methods, and document\n" +
			"types. This is the data the official MCP's units/reference tools serve; Alegra\n" +
			"exposes no public REST endpoint for it, so the CLI embeds it (generated from\n" +
			"Alegra's published country parameter pages).\n\n" +
			"Run without an argument to list the categories for your country; pass a category\n" +
			"key to list its values. The country is auto-detected from your account; override\n" +
			"it with --country (works offline, no login required).",
		Example: "  alegra catalog                       # list categories for your country\n" +
			"  alegra catalog units                 # units of measure\n" +
			"  alegra catalog identification-types  # valid contact ID types\n" +
			"  alegra catalog units --country mexico -o json",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cat, err := catalog.Load(resolveCountry(country))
			if err != nil {
				return err
			}
			if len(args) == 0 {
				rows := make([]catalogCategoryRow, 0, len(cat.Categories))
				for _, c := range cat.Categories {
					rows = append(rows, catalogCategoryRow{Category: c.Key, Title: c.Title, Entries: len(c.Entries)})
				}
				sort.Slice(rows, func(i, j int) bool { return rows[i].Category < rows[j].Category })
				return render(cmd, rows, []string{"category", "title", "entries"})
			}
			c, ok := cat.Category(args[0])
			if !ok {
				return fmt.Errorf("unknown category %q for %s; available: %v", args[0], cat.Label, cat.CategoryKeys())
			}
			return render(cmd, c.Entries, []string{"code", "name"})
		},
	}
	cmd.Flags().StringVar(&country, "country", "", "Country to look up (default: auto-detected from the account; e.g. colombia, mexico, costaRica, peru)")
	_ = cmd.RegisterFlagCompletionFunc("country", fixedCompleter(catalog.Available()...))
	cmd.AddCommand(newCatalogSyncSATCmd(), newCatalogProductKeysCmd())
	rootCmd.AddCommand(cmd)
}

// catalogsDir is the shared (cross-profile) local catalog cache, next to the
// config file: the SAT catalog is global government data, not account data.
func catalogsDir() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(cfg.Path()), "catalogs"), nil
}

func newCatalogSyncSATCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-sat",
		Short: "Download the SAT product-keys catalog (México) to the local cache",
		Long: "Download México's SAT c_ClaveProdServ catalog (~52k product/service keys,\n" +
			"~7MB) into the shared local cache. Alegra exposes no endpoint for it, so the\n" +
			"CLI sources it from the SAT's published data (phpcfdi/resources-sat-catalogs\n" +
			"mirror). Needed once; re-run to pick up new SAT catalog versions. No Alegra\n" +
			"credentials required.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			dir, err := catalogsDir()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Downloading SAT product-keys catalog (~7MB)…")
			cat, err := catalog.SyncSAT(cmd.Context(), dir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Synced %d product keys (catalog version %s) → %s\n",
				len(cat.Entries), cat.Version, catalog.SATPath(dir))
			return nil
		},
	}
}

func newCatalogProductKeysCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:     "product-keys [query]",
		Aliases: []string{"product-key", "claves-producto"},
		Short:   "Search SAT product/service keys (México, claveProdServ)",
		Long: "Search México's SAT c_ClaveProdServ catalog offline. Matching is case- and\n" +
			"accent-insensitive across the key, its description, and the SAT's published\n" +
			"similar-names list. Requires a one-time `alegra catalog sync-sat`.",
		Example: "  alegra catalog product-keys refrigerador\n" +
			"  alegra catalog product-keys 10101506 -o json\n" +
			"  alegra catalog product-keys \"servicios de programacion\" --limit 5",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := catalogsDir()
			if err != nil {
				return err
			}
			cat, err := catalog.LoadSAT(dir)
			if err != nil {
				return err
			}
			query := ""
			if len(args) == 1 {
				query = args[0]
			}
			entries := catalog.SearchSAT(cat, query, limit)
			if len(entries) == 0 {
				return fmt.Errorf("no product keys match %q (catalog version %s, %d entries)", query, cat.Version, len(cat.Entries))
			}
			return render(cmd, entries, []string{"code", "name"})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum results to show (0 = no limit)")
	return cmd
}
