package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.GlobalInvoice]{
		Use:         "global-invoices",
		Aliases:     []string{"global-invoice"},
		Short:       "Manage global invoices (facturas globales)",
		Long:        "Manage Alegra global invoices, which group sale tickets into a single CFDI for the general public.",
		New:         func(c *api.Client) *api.Resource[api.GlobalInvoice] { return c.GlobalInvoices() },
		Columns:     []string{"id", "date", "status", "total"},
		OrderFields: []string{"id", "date", "status"},
		ListFilters: []listFilter{
			{Flag: "status", Query: "status", Usage: "Filter by status: open, draft, void"},
			{Flag: "date", Query: "date", Usage: "Filter by date (YYYY-MM-DD)"},
		},
	})
}
