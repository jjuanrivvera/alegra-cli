package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Term]{
		Use:         "terms",
		Aliases:     []string{"term"},
		Short:       "Manage payment terms (términos de pago)",
		Long:        "Manage payment terms (términos de pago) — named due-date rules (e.g. net 30) applied to invoices and bills.",
		New:         func(c *api.Client) *api.Resource[api.Term] { return c.Terms() },
		Columns:     []string{"id", "name", "days", "status"},
		OrderFields: []string{"id", "name", "days"},
		ListFilters: []listFilter{
			{Flag: "fields", Query: "fields", Usage: "Extra fields to include (e.g. deletable)"},
		},
	})
}
