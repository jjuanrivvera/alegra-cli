package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Estimate]{
		Use:         "estimates",
		Aliases:     []string{"estimate", "quote"},
		Short:       "Manage estimates (cotizaciones / sales quotes)",
		New:         func(c *api.Client) *api.Resource[api.Estimate] { return c.Estimates() },
		Columns:     []string{"id", "date", "status", "total"},
		OrderFields: []string{"id", "name", "date", "dueDate"},
		ListFilters: append([]listFilter{
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "client-name", Query: "client_name", Usage: "Filter by client name"},
			{Flag: "item-id", Query: "item_id", Usage: "Filter by item ID"},
			{Flag: "number", Query: "number", Usage: "Filter by estimate number (prefix and number)"},
			{Flag: "date", Query: "date", Usage: "Filter by exact date (YYYY-MM-DD)"},
		}, dateRangeFilters()...),
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Estimate]) {
			parent.AddCommand(NewActionCmd(sp, "email", "email", "Email an estimate", true))
		},
	})
}
