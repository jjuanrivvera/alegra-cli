package commands

import "github.com/jjuanrivvera/alegra-cli/internal/api"

func init() {
	registerResource(resourceSpec[api.Item]{
		Use:         "items",
		Aliases:     []string{"item"},
		Short:       "Manage items (products and services)",
		New:         func(c *api.Client) *api.Resource[api.Item] { return c.Items() },
		Columns:     []string{"id", "name", "reference", "price", "status"},
		OrderFields: []string{"id", "name", "reference", "price"},
		ListFilters: []listFilter{
			{Flag: "name", Query: "name", Usage: "Filter by name"},
		},
	})
}
