package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Payment]{
		Use:         "payments",
		Aliases:     []string{"payment"},
		Short:       "Manage payments (incomes and expenses)",
		New:         func(c *api.Client) *api.Resource[api.Payment] { return c.Payments() },
		Columns:     []string{"id", "date", "amount", "type", "status"},
		OrderFields: []string{"id", "date"},
		ListFilters: append([]listFilter{
			{Flag: "type", Query: "type", Usage: "in or out", Values: []string{"in", "out"}},
			{Flag: "status", Query: "status", Usage: "Filter by status"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
		}, dateRangeFilters()...),
		Extra: func(parent *cobra.Command, sp resourceSpec[api.Payment]) {
			parent.AddCommand(NewActionCmd(sp, "void", "void", "Void a payment"))
			parent.AddCommand(NewActionCmd(sp, "open", "open", "Revert a voided payment"))
			parent.AddCommand(NewActionCmd(sp, "stamp", "stamp", "Emit electronic payment receipt (REP)"))
		},
	})
}
