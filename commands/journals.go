package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Journal]{
		Use:         "journals",
		Aliases:     []string{"journal"},
		Short:       "Manage accounting journal entries (comprobantes contables)",
		Long:        "Manage manual accounting journal entries (comprobantes contables) — direct debit/credit postings to ledger accounts. Use `alegra journals balance` for balances grouped by period.",
		New:         func(c *api.Client) *api.Resource[api.Journal] { return c.Journals() },
		Columns:     []string{"id", "date", "number", "status"},
		OrderFields: []string{"date", "name", "reference", "observations"},
		ListFilters: []listFilter{
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "client-name", Query: "client_name", Usage: "Filter by client name"},
			{Flag: "reference", Query: "reference", Usage: "Filter by reference"},
			{Flag: "observations", Query: "observations", Usage: "Filter by observations"},
			{Flag: "number", Query: "number", Usage: "Filter by numbering"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Journal]) {
			parent.AddCommand(&cobra.Command{
				Use:   "balance",
				Short: "Retrieve journal balances grouped by month or day",
				Args:  cobra.NoArgs,
				RunE: func(cmd *cobra.Command, _ []string) error {
					client, err := getAPIClient()
					if err != nil {
						return err
					}
					var out map[string]any
					if err := client.GetInto(cmd.Context(), "journals/entries/graph", nil, &out); err != nil {
						return err
					}
					return render(cmd, out, nil)
				},
			})
		},
	})
}
