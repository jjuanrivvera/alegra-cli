package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.RecurringInvoice]{
		Use:     "recurring-invoices",
		Aliases: []string{"recurring-invoice"},
		Short:   "Manage recurring invoices (facturas recurrentes)",
		Long:    "Manage recurring invoices: templates that automatically generate sales invoices on a schedule.",
		New:     func(c *api.Client) *api.Resource[api.RecurringInvoice] { return c.RecurringInvoices() },
		Columns: []string{"id", "status", "total"},
		OrderFields: []string{
			"name", "startDate", "endDate", "repeatEvery", "term",
		},
		ListFilters: []listFilter{
			{Flag: "start-date", Query: "startDate", Usage: "Filter by recurring invoice start date (YYYY-MM-DD)"},
			{Flag: "end-date", Query: "endDate", Usage: "Filter by recurring invoice end date (YYYY-MM-DD)"},
			{Flag: "repeat-every", Query: "repeatEvery", Usage: "Filter by recurrence interval"},
			{Flag: "term", Query: "term", Usage: "Filter by payment term"},
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
			{Flag: "name", Query: "name", Usage: "Filter by client name"},
		},
	})
}
