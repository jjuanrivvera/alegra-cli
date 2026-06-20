package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.RecurringPayment]{
		Use:     "recurring-payments",
		Aliases: []string{"recurring-payment"},
		Short:   "View recurring payments (pagos recurrentes)",
		Long:    "View Alegra recurring payments, which are automatically registered against a bank account on a recurring schedule. This resource is read-only.",
		New:     func(c *api.Client) *api.Resource[api.RecurringPayment] { return c.RecurringPayments() },
		Columns: []string{"id", "status"},
		// Recurring payments are charged against a bank account, not a client, and
		// the record carries startDate/endDate (no "date"/"type" fields), so the
		// order fields and the copied "client-id" filter were both inapplicable.
		OrderFields: []string{"id", "number", "startDate", "endDate"},
		NoCreate:    true,
		NoUpdate:    true,
		NoDelete:    true,
	})
}
