package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.RecurringPayment]{
		Use:         "recurring-payments",
		Aliases:     []string{"recurring-payment"},
		Short:       "View recurring payments (pagos recurrentes)",
		Long:        "View Alegra recurring payments, which are automatically registered against a bank account on a recurring schedule. This resource is read-only.",
		New:         func(c *api.Client) *api.Resource[api.RecurringPayment] { return c.RecurringPayments() },
		Columns:     []string{"id", "status"},
		OrderFields: []string{"id", "number", "date", "type"},
		ListFilters: []listFilter{
			{Flag: "client-id", Query: "client_id", Usage: "Filter by client ID"},
		},
		NoCreate: true,
		NoUpdate: true,
		NoDelete: true,
	})
}
