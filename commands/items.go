package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Item]{
		Use:         "items",
		Aliases:     []string{"item"},
		Short:       "Manage items (products and services)",
		Long:        "Manage items (products and services) — your catalog. Items carry price, taxes, reference codes, and inventory settings, and become the line items on invoices and bills.",
		New:         func(c *api.Client) *api.Resource[api.Item] { return c.Items() },
		Columns:     []string{"id", "name", "reference", "price", "status"},
		OrderFields: []string{"id", "name", "reference", "price"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter by name"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Item]) {
			parent.AddCommand(readOnlyHints(itemStockCmd(sp)))
		},
	})
}

// itemStockCmd shows an item's stock broken down by warehouse, read from the
// item's inventory (GET /items/{id}). Use --date for a historical snapshot.
func itemStockCmd(sp resourceSpec[api.Item]) *cobra.Command {
	var date string
	cmd := &cobra.Command{
		Use:               "stock <id>",
		Short:             "Show an item's stock, broken down by warehouse",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: resourceIDCompleter(sp),
		Example:           "  alegra items stock <id>\n  alegra items stock <id> --date 2026-01-31 -o json",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient(cmd)
			if err != nil {
				return err
			}
			q := url.Values{}
			if date != "" {
				q.Set("documentDate", date)
			}
			var item api.Item
			if err := client.GetInto(cmd.Context(), "items/"+url.PathEscape(args[0]), q, &item); err != nil {
				return err
			}
			if item.Inventory == nil {
				if !flagDryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "Item %s is not inventariable (no stock tracked)\n", args[0])
				}
				return nil
			}
			if len(item.Inventory.Warehouses) > 0 {
				return render(cmd, item.Inventory.Warehouses, []string{"id", "name", "availableQuantity", "initialQuantity"})
			}
			return render(cmd, item.Inventory, nil)
		},
	}
	cmd.Flags().StringVar(&date, "date", "", "Historical stock as of this date (YYYY-MM-DD)")
	return cmd
}
