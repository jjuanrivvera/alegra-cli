package commands

import (
	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	registerResource(resourceSpec[api.IncomeDebitNote]{
		Use:         "income-debit-notes",
		Aliases:     []string{"income-debit-note"},
		Short:       "Manage customer debit notes",
		Long:        "Manage customer debit notes (notas débito de cliente) — adjustments that increase what a customer owes, issued against income documents. For supplier-side notes use `alegra debit-notes`.",
		New:         func(c *api.Client) *api.Resource[api.IncomeDebitNote] { return c.IncomeDebitNotes() },
		Columns:     []string{"id", "date", "status", "total"},
		OrderFields: []string{"id", "date", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: open, closed, void"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
		},
	})
}
