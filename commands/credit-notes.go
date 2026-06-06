package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.CreditNote]{
		Use:         "credit-notes",
		Aliases:     []string{"credit-note"},
		Short:       "Manage credit notes",
		New:         func(c *api.Client) *api.Resource[api.CreditNote] { return c.CreditNotes() },
		Columns:     []string{"id", "date", "status", "total"},
		OrderFields: []string{"id", "date", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: open, closed, void"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.CreditNote]) {
			parent.AddCommand(NewActionCmd(sp, "email", "email", "Email a credit note"))
		},
	})
}
