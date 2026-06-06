package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Bill]{
		Use:         "bills",
		Aliases:     []string{"bill"},
		Short:       "Manage provider bills (facturas de proveedor)",
		New:         func(c *api.Client) *api.Resource[api.Bill] { return c.Bills() },
		Columns:     []string{"id", "date", "dueDate", "status", "total", "balance"},
		OrderFields: []string{"id", "date", "dueDate", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status"},
			{Flag: "provider-name", Query: "provider_name", Usage: "Filter by provider name"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Bill]) {
			parent.AddCommand(NewActionCmd(sp, "close", "close", "Close a bill with pending balance"))
			parent.AddCommand(NewActionCmd(sp, "comments", "comments", "Add a comment to a bill"))
			parent.AddCommand(NewCollectionActionCmd(sp, "import-by-cufe", "import-by-cufe", "Import a bill by CUFE (Colombia)"))
		},
	})
}
