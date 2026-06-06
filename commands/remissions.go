package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Remission]{
		Use:         "remissions",
		Aliases:     []string{"remission"},
		Short:       "Manage remissions (delivery notes)",
		New:         func(c *api.Client) *api.Resource[api.Remission] { return c.Remissions() },
		Columns:     []string{"id", "date", "status"},
		OrderFields: []string{"id", "name", "date", "dueDate"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: open, void, closed"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "client-name", Query: "client_name", Usage: "Filter by client name"},
			{Flag: "item-id", Query: "item_id", Usage: "Filter by item ID"},
			{Flag: "number", Query: "number", Usage: "Filter by remission number (prefix and number)"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
			{Flag: "due-date", Query: "dueDate", Usage: "Filter by due date (YYYY-MM-DD)"},
		},
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Remission]) {
			parent.AddCommand(NewActionCmd(sp, "void", "void", "Void a remission"))
			parent.AddCommand(NewActionCmd(sp, "open", "open", "Open a remission"))
		},
	})
}
