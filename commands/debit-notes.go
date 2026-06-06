package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.DebitNote]{
		Use:         "debit-notes",
		Aliases:     []string{"debit-note"},
		Short:       "Manage debit notes",
		New:         func(c *api.Client) *api.Resource[api.DebitNote] { return c.DebitNotes() },
		Columns:     []string{"id", "date", "status", "total"},
		OrderFields: []string{"id", "date", "status"},
		ListFilters: []listFilter{
			{Flag: "number", Query: "number", Usage: "Filter by document number (prefix and/or number)"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "item-id", Query: "item_id", Usage: "Filter by item ID"},
			{Flag: "provider-name", Query: "provider_name", Usage: "Filter by provider name"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.DebitNote]) {
			parent.AddCommand(NewActionCmd(sp, "comments", "comments", "Add a comment to a debit note", true))
		},
	})
}
