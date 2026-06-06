package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Category]{
		Use:         "categories",
		Aliases:     []string{"account"},
		Short:       "Manage chart-of-accounts accounts (cuentas contables)",
		Long:        "Manage chart-of-accounts accounts (cuentas contables). These are accounting ledger accounts. There is no delete endpoint.",
		New:         func(c *api.Client) *api.Resource[api.Category] { return c.Categories() },
		Columns:     []string{"id", "name", "code", "type", "status"},
		OrderFields: []string{"id", "name", "code"},
		NoDelete:    true,
		ListFilters: []listFilter{
			{Flag: "type", Query: "type", Usage: "Filter by type: income, expense, asset, liability, equity, cost, productionCost, order"},
			{Flag: "status", Query: "status", Usage: "Filter by status: active, inactive, deleted"},
			{Flag: "name", Query: "name", Usage: "Filter by account name"},
			{Flag: "format", Query: "format", Usage: "Response format (e.g. tree) to show account hierarchy"},
		},
		Extra: func(parent *cobra.Command, _ resourceSpec[api.Category]) {
			parent.AddCommand(newCategoriesSettingsCmd())
			parent.AddCommand(newCategoriesSetSettingsCmd())
		},
	})
}

// newCategoriesSettingsCmd builds `categories settings` → GET /categories/settings.
func newCategoriesSettingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Show chart-of-accounts settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			var out any
			if err := client.GetInto(cmd.Context(), "categories/settings", nil, &out); err != nil {
				return err
			}
			return render(cmd, out, nil)
		},
	}
}

// newCategoriesSetSettingsCmd builds `categories set-settings` → PUT /categories/settings.
func newCategoriesSetSettingsCmd() *cobra.Command {
	var bf bodyFlags
	cmd := &cobra.Command{
		Use:   "set-settings",
		Short: "Update chart-of-accounts settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := bf.build()
			if err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			var out any
			if err := client.PutInto(cmd.Context(), "categories/settings", body, &out); err != nil {
				return err
			}
			if flagDryRun {
				return nil
			}
			if out == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Settings updated.")
				return nil
			}
			return render(cmd, out, nil)
		},
	}
	addBodyFlags(cmd, &bf)
	return cmd
}
