package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Invoice]{
		Use:         "invoices",
		Aliases:     []string{"invoice"},
		Short:       "Manage sales invoices",
		New:         func(c *api.Client) *api.Resource[api.Invoice] { return c.Invoices() },
		Columns:     []string{"id", "date", "dueDate", "status", "total", "balance"},
		OrderFields: []string{"id", "name", "date", "dueDate", "status"},
		ListFilters: append([]listFilter{
			{Flag: "status", Query: "status", Usage: "open,closed,draft,void"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "date", Query: "date", Usage: "Filter by exact date YYYY-MM-DD"},
		}, dateRangeFilters()...),
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Invoice]) {
			parent.AddCommand(newInvoicesEmitCmd(sp))
			parent.AddCommand(NewActionCmd(sp, "void", "void", "Void an invoice"))
			parent.AddCommand(NewActionCmd(sp, "open", "open", "Revert a voided invoice"))
			parent.AddCommand(NewActionCmd(sp, "email", "email", "Email an invoice", true))
			parent.AddCommand(NewCollectionActionCmd(sp, "stamp", "stamp", "Stamp up to 10 invoices (low-level; prefer `emit`)"))
			parent.AddCommand(NewCollectionActionCmd(sp, "preview", "preview", "Generate a preview PDF URL"))
		},
	})
}
