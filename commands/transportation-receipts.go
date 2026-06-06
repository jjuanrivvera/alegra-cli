package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.TransportationReceipt]{
		Use:         "transportation-receipts",
		Aliases:     []string{"transportation-receipt"},
		Short:       "Manage transportation receipts (documentos de traslado)",
		New:         func(c *api.Client) *api.Resource[api.TransportationReceipt] { return c.TransportationReceipts() },
		Columns:     []string{"id", "date", "status"},
		OrderFields: []string{"id", "date", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: open, closed, draft, void"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.TransportationReceipt]) {
			parent.AddCommand(NewActionCmd(sp, "void", "void", "Void a receipt"))
			parent.AddCommand(NewActionCmd(sp, "open", "open", "Open a receipt"))
			parent.AddCommand(NewActionCmd(sp, "email", "email", "Email a receipt", true))
			parent.AddCommand(NewCollectionActionCmd(sp, "preview", "preview", "Generate a preview PDF URL"))
		},
	})
}
