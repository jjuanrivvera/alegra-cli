package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.Payment]{
		Use:     "payments",
		Aliases: []string{"payment"},
		Short:   "Manage payments (incomes and expenses)",
		Long: "Manage payments. type \"in\" is money received (income, against invoices); " +
			"type \"out\" is money paid (expense, against bills). Allocate a payment to documents " +
			"via the invoices[]/bills[] arrays in the body. In Costa Rica/Mexico a payment may need " +
			"electronic stamping — see the `stamp` subcommand (REP / complemento de pago).",
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
