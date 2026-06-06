package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Reconciliation]{
		Use:         "reconciliations",
		Aliases:     []string{"reconciliation"},
		Short:       "Manage bank reconciliations",
		Long:        "Manage bank reconciliations. Create performs a create-or-update (upsert); there is no separate update endpoint.",
		New:         func(c *api.Client) *api.Resource[api.Reconciliation] { return c.Reconciliations() },
		Columns:     []string{"id", "status"},
		OrderFields: []string{"id", "date"},
		NoUpdate:    true,
		ListFilters: []listFilter{
			{Flag: "account-id", Query: "account_id", Usage: "Filter by bank account ID"},
			{Flag: "fields", Query: "fields", Usage: "Extra fields: deletable, balance, or simple"},
		},
	})
}
