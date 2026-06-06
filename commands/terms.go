package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Term]{
		Use:         "terms",
		Aliases:     []string{"term"},
		Short:       "Manage payment terms (términos de pago)",
		New:         func(c *api.Client) *api.Resource[api.Term] { return c.Terms() },
		Columns:     []string{"id", "name", "days", "status"},
		OrderFields: []string{"id", "name", "days"},
		ListFilters: []listFilter{
			{Flag: "fields", Query: "fields", Usage: "Extra fields to include (e.g. deletable)"},
		},
	})
}
