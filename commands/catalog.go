package commands

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/catalog"
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
	rootCmd.AddCommand(cmd)
}
